package server

import (
	"context"
	"fmt"
	"github.com/msalvati1997/b1multicasting/pkg/reg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"strings"
	"sync"
)

type server struct {
	proto.UnimplementedRegistryServer
}

var (
	groups   map[string]*Mgroup
	Mugroups sync.RWMutex
)

type Mgroup struct {
	mu        sync.RWMutex
	groupInfo *proto.MGroup
}

func init() {
	groups = make(map[string]*Mgroup)
}

// Registration registers the calling node to a multicast group
func (s *server) Register(ctx context.Context, in *proto.Rinfo) (*proto.Ranswer, error) {
	source, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	src := source.Addr.String()
	srcAddr := src[:strings.LastIndexByte(src, ':')]
	srcAddr = fmt.Sprintf("%s:%d", srcAddr, in.ClientPort)
	log.Println("Registration of the group ", in.MulticastId, "with client", srcAddr)
	multicastId := in.MulticastId
	Mugroups.Lock()
	defer Mugroups.Unlock()
	// Checking if the group already exists
	multicastGroup, exists := groups[multicastId]

	if exists {
		if multicastGroup.groupInfo.Status != proto.Status_OPENING {
			return nil, status.Errorf(codes.PermissionDenied, "ServiceStarted")
		}
		if _, ok := multicastGroup.groupInfo.Members[srcAddr]; ok {
			return nil, status.Errorf(codes.AlreadyExists, "Already Registered")
		}
	} else {
		log.Println("Start registering ..")
		// Creating the group
		multicastGroup = &Mgroup{groupInfo: &proto.MGroup{
			MulticastId: multicastId,
			Status:      proto.Status_OPENING,
			Members:     make(map[string]*proto.MemberInfo),
		}}
		groups[multicastId] = multicastGroup
	}
	// Registering node to the group
	//Creating new MemberInfo
	memberInfo := new(proto.MemberInfo)
	memberInfo.Id = srcAddr
	memberInfo.Address = srcAddr
	memberInfo.Ready = false
	//Adding the MemberInfo to the multicastGroup Map
	multicastGroup.groupInfo.Members[srcAddr] = memberInfo

	return &proto.Ranswer{
		ClientId:  srcAddr,
		GroupInfo: multicastGroup.groupInfo,
	}, nil
}

// Start enables multicast in the group.
// All processes belonging to the group can initializes the necessary structure and
// then they must communicate they are ready using Ready function.
func (s *server) startGroup(ctx context.Context, in *proto.RequestData) (*proto.MGroup, error) {
	_, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	multicastId := in.MulticastId
	clientId := in.MId
	Mugroups.RLock()
	defer Mugroups.RUnlock()
	// Checking if the group exists
	mGroup, ok := groups[multicastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}
	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()
	// Checking if the group can be started
	if mGroup.groupInfo.Status != proto.Status_OPENING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}
	// Checking if the process belong to the group
	if _, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}
	// Updating group status
	mGroup.groupInfo.Status = proto.Status_STARTING
	return mGroup.groupInfo, nil
}

// Ready stats that the calling node has correctly initialized the required structures for multicast.
// Nodes can multicast messages when all members have called the Ready function
func (s *server) Ready(ctx context.Context, in *proto.RequestData) (*proto.MGroup, error) {

	_, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	multicastId := in.MulticastId
	clientId := in.MId
	Mugroups.RLock()
	defer Mugroups.RUnlock()

	// Checking if the group exists
	mGroup, ok := groups[multicastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()

	// Checking if the group is in the correct status
	if mGroup.groupInfo.Status != proto.Status_STARTING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	var member *proto.MemberInfo

	// Checking if the process belongs to the group
	if member, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	if !member.Ready {
		member.Ready = true
		mGroup.groupInfo.ReadyMembers++
	}

	if mGroup.groupInfo.ReadyMembers == uint64(len(mGroup.groupInfo.Members)) {
		// All members are ready, they can multicast messages
		mGroup.groupInfo.Status = proto.Status_ACTIVE
	}
	return mGroup.groupInfo, nil
}

// CloseGroup closes the group. The process cannot longer use the group to multicast messages
func (s *server) CloseGroup(ctx context.Context, in *proto.RequestData) (*proto.MGroup, error) {
	_, ok := peer.FromContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}

	multicastId := in.MulticastId
	clientId := in.MId

	Mugroups.RLock()
	defer Mugroups.RUnlock()

	// Checking if the group exists
	mGroup, ok := groups[multicastId]

	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.Lock()
	defer mGroup.mu.Unlock()

	// Checking if the group can be closed
	if mGroup.groupInfo.Status != proto.Status_ACTIVE && mGroup.groupInfo.Status != proto.Status_CLOSING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	// Checking if the node belongs to the group
	member, ok := mGroup.groupInfo.Members[clientId]

	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = proto.Status_CLOSING

	if member.Ready {
		member.Ready = false
		mGroup.groupInfo.ReadyMembers--
		if mGroup.groupInfo.ReadyMembers == 0 {
			// Closing definitely the group when all members have called CloseGroup function
			mGroup.groupInfo.Status = proto.Status_CLOSED
		}
	}
	return mGroup.groupInfo, nil
}

func Registration(s grpc.ServiceRegistrar) (err error) {
	proto.RegisterRegistryServer(s, &server{})
	return
}

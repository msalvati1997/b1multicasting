package server

import (
	"fmt"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RegistryServer struct {
	protoregistry.UnimplementedRegistryServer
}

var (
	groups   map[string]*Mgroup
	Mugroups sync.RWMutex
)

type Mgroup struct {
	mu        sync.RWMutex
	groupInfo *protoregistry.MGroup
}

func init() {
	groups = make(map[string]*Mgroup)
}

// Registration registers the calling node to a multicast group
func (s *RegistryServer) Register(ctx context.Context, in *protoregistry.Rinfo) (*protoregistry.Ranswer, error) {
	timeout := time.Second
	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	log.Println("Start registering..")
	source, ok := peer.FromContext(c)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}
	src := source.Addr.String()
	srcAddr := src[:strings.LastIndexByte(src, ':')]
	id := strings.Split(src[:strings.LastIndexByte(src, ':')], ".")
	id1 := strings.Join(id, "")
	ids, _ := strconv.Atoi(id1)
	utils.Myid = ids
	utils.MystringId = id1
	log.Println("Trying to convert...", ids)
	srcAddr = fmt.Sprintf("%s:%d", srcAddr, in.ClientPort)
	utils.MyAdress = srcAddr
	log.Println("Registration of the group ", in.MulticastId, "with client", srcAddr)
	multicastId := in.MulticastId
	Mugroups.Lock()
	defer Mugroups.Unlock()
	// Checking if the group already exists
	multicastGroup, exists := groups[multicastId]

	if exists {
		if multicastGroup.groupInfo.Status != protoregistry.Status_OPENING {
			return nil, status.Errorf(codes.PermissionDenied, "ServiceStarted")
		}
		if _, ok := multicastGroup.groupInfo.Members[srcAddr]; ok {
			return nil, status.Errorf(codes.AlreadyExists, "Already Registered")
		}
	} else {
		log.Println("Start registering ..")
		// Creating the group
		multicastGroup = &Mgroup{groupInfo: &protoregistry.MGroup{
			MulticastId: multicastId,
			Status:      protoregistry.Status_OPENING,
			Members:     make(map[string]*protoregistry.MemberInfo),
		}}
		groups[multicastId] = multicastGroup
	}
	// Registering node to the group
	//Creating new MemberInfo
	memberInfo := new(protoregistry.MemberInfo)
	memberInfo.Id = srcAddr
	memberInfo.Address = srcAddr
	memberInfo.Ready = false
	//Adding the MemberInfo to the multicastGroup Map
	multicastGroup.groupInfo.Members[srcAddr] = memberInfo

	return &protoregistry.Ranswer{
		ClientId:  srcAddr,
		GroupInfo: multicastGroup.groupInfo,
	}, nil
}

// Start enables multicast in the group.
// All processes belonging to the group can initializes the necessary structure and
// then they must communicate they are ready using Ready function.
func (s *RegistryServer) StartGroup(ctx context.Context, in *protoregistry.RequestData) (*protoregistry.MGroup, error) {
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
	if mGroup.groupInfo.Status != protoregistry.Status_OPENING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}
	// Checking if the process belong to the group
	if _, ok = mGroup.groupInfo.Members[clientId]; !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}
	// Updating group status
	mGroup.groupInfo.Status = protoregistry.Status_STARTING
	return mGroup.groupInfo, nil
}

// Ready stats that the calling node has correctly initialized the required structures for multicast.
// Nodes can multicast messages when all members have called the Ready function
func (s *RegistryServer) Ready(ctx context.Context, in *protoregistry.RequestData) (*protoregistry.MGroup, error) {

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
	if mGroup.groupInfo.Status != protoregistry.Status_STARTING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	var member *protoregistry.MemberInfo

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
		mGroup.groupInfo.Status = protoregistry.Status_ACTIVE
	}
	return mGroup.groupInfo, nil
}

// CloseGroup closes the group. The process cannot longer use the group to multicast messages
func (s *RegistryServer) CloseGroup(ctx context.Context, in *protoregistry.RequestData) (*protoregistry.MGroup, error) {

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
	if mGroup.groupInfo.Status != protoregistry.Status_ACTIVE && mGroup.groupInfo.Status != protoregistry.Status_CLOSING {
		return nil, status.Errorf(codes.Canceled, "Cannot start group")
	}

	// Checking if the node belongs to the group
	member, ok := mGroup.groupInfo.Members[clientId]

	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Not registered to the group")
	}

	mGroup.groupInfo.Status = protoregistry.Status_CLOSING

	if member.Ready {
		member.Ready = false
		mGroup.groupInfo.ReadyMembers--
		if mGroup.groupInfo.ReadyMembers == 0 {
			// Closing definitely the group when all members have called CloseGroup function
			mGroup.groupInfo.Status = protoregistry.Status_CLOSED
		}
	}
	return mGroup.groupInfo, nil
}

func Registration(s grpc.ServiceRegistrar) (err error) {
	protoregistry.RegisterRegistryServer(s, &RegistryServer{})
	return
}

// GetStatus returns infos about the group associated to multicastId
func (s *RegistryServer) GetStatus(ctx context.Context, in *protoregistry.MulticastId) (*protoregistry.MGroup, error) {
	_, ok := peer.FromContext(ctx)

	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing source address")
	}

	Mugroups.RLock()
	defer Mugroups.RUnlock()

	mGroup, ok := groups[in.MulticastId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "Group not found")
	}

	mGroup.mu.RLock()
	defer mGroup.mu.RUnlock()

	return mGroup.groupInfo, nil

}

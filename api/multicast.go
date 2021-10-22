package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/msalvati1997/b1multicasting/pkg/reg/proto"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	Registryclient  proto.RegistryClient
	GMu             sync.RWMutex
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
	groupsName      []MulticastId
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	ClientId  string `json:"client_id"`
	Group     *MulticastInfo
	GroupMu   sync.RWMutex `json:"-"`
	Messages  []Message
	MessageMu sync.RWMutex `json:"-"`
}

type Message struct {
	MessageHeader map[string]string `json:"MessageHeader"`
	Payload       []byte            `json:"Payload"`
}

type MulticastInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
}

// ErrorResponse defines an error response returned upon any request failure.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ErrorResponse defines an error response returned upon any request failure.
func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	bz, _ := json.Marshal(ErrorResponse{Error: err.Error()})
	_, _ = w.Write(bz)
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

type MulticastId struct {
	MulticastId string `json:"multicast_id"`
}

// GetGroups godoc
// @Summary Get details of all groups
// @Description Get details of all groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {object} MulticastGroup
// @Router /groups [get]
func GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(MulticastGroups)
	if err != nil {
		return
	}
}

// CreateGroup godoc
// @Summary Create a group from an id
// @Description Create a group from an id
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {object} MulticastGroup
// @Router /groups [put]
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var multicastId MulticastId
	err := json.NewDecoder(r.Body).Decode(&multicastId)
	if err != nil {
		return
	}
	groupsName = append(groupsName, multicastId)
	log.Println("Start creating group with ", multicastId.MulticastId, " at grpcPort", uint32(GrpcPort))
	w.Header().Set("Content-Type", "application/json")
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[multicastId.MulticastId]

	if ok {
		log.Println("The group already exist")
	} else {
		log.Println("The group doesn't exist before")
	}
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()
	r = r.WithContext(ctx)

	register, err := Registryclient.Register(r.Context(), &proto.Rinfo{
		MulticastId: multicastId.MulticastId,
		ClientPort:  uint32(GrpcPort),
	})
	if err != nil {
		log.Println("Problem in regitering member to group", err.Error())
	} else {
		log.Println("ok in registering")
	}
	members := make(map[string]Member, 0)

	for memberId, member := range register.GroupInfo.Members {
		members[memberId] = Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &MulticastGroup{
		ClientId: register.ClientId,
		Group: &MulticastInfo{
			MulticastId:      register.GroupInfo.MulticastId,
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(register.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]Message, 0),
	}

	MulticastGroups[register.GroupInfo.MulticastId] = group
	go InitGroup(register.GroupInfo, group, len(register.GroupInfo.Members) == 1)
	log.Println("Group created")

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(MulticastGroups)
	if err != nil {
		return
	}

}

func InitGroup(info *proto.MGroup, group *MulticastGroup, b bool) {
	// Waiting that the group is ready
	update(info, group)
	groupInfo, _ := StatusChange(info, group, proto.Status_OPENING)

	// Communicating to the registry that the node is ready
	groupInfo, _ = Registryclient.Ready(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.ClientId,
	})

	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, proto.Status_STARTING)

}

func StatusChange(groupInfo *proto.MGroup, multicastGroup *MulticastGroup, status proto.Status) (*proto.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = Registryclient.GetStatus(context.Background(), &proto.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}

		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == proto.Status_CLOSED || groupInfo.Status == proto.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

func update(groupInfo *proto.MGroup, multicastGroup *MulticastGroup) {
	multicastGroup.GroupMu.Lock()
	defer multicastGroup.GroupMu.Unlock()

	multicastGroup.Group.Status = proto.Status_name[int32(groupInfo.Status)]

	for clientId, member := range groupInfo.Members {
		m, ok := multicastGroup.Group.Members[clientId]

		if !ok {
			m = Member{
				MemberId: member.Id,
				Address:  member.Address,
				Ready:    member.Ready,
			}

			multicastGroup.Group.Members[clientId] = m
		}

		if m.Ready != member.Ready {
			m.Ready = member.Ready
			multicastGroup.Group.Members[clientId] = m
		}

	}
}

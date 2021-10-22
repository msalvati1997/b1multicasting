package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/reg"
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
	groupsName      []MulticastReq
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

type MulticastReq struct {
	MulticastId   string `json:"multicast_id"`
	MulticastType proto.MulticastType
}

func GetGroups(g *gin.Context) {
	groups := make([]*MulticastInfo, 0)

	GMu.RLock()
	defer GMu.RUnlock()

	for _, group := range MulticastGroups {
		group.GroupMu.RLock()
		groups = append(groups, group.Group)
		group.GroupMu.RUnlock()
	}

	g.IndentedJSON(http.StatusOK, groups)
}

// CreateGroup godoc
// @Summary Create Multicast Group
// @Description Create Multicast Group
// @Tags groups
// @Accept  json
// @Produce  json
// @Param post body MulticastReq
// @Success 201 {object} MulticastGroups
//     Responses:
//       201: body:PositionResponseBody
// @Router /groups [post]
func CreateGroup(g *gin.Context) {
	var multicastId MulticastReq
	err := json.NewDecoder(g.Request.Body).Decode(&multicastId)
	if err != nil {
		g.JSON(http.StatusBadRequest, "Error in body request")
	}
	groupsName = append(groupsName, multicastId)

	log.Println("Start creating group with ", multicastId.MulticastId, " and multicast type ", multicastId.MulticastType, "at grpcPort", uint32(GrpcPort))

	multicastType, ok := reg.MulticastType[multicastId.MulticastType.String()]
	if !ok {
		g.IndentedJSON(http.StatusBadRequest, gin.H{"error": "multicast_type not supported", "supported": reg.MulticastTypes})
		return
	}
	GMu.RLock()
	defer GMu.RUnlock()

	id := multicastId.MulticastId
	group, ok := MulticastGroups[id]

	if ok {
		g.IndentedJSON(http.StatusConflict, gin.H{"error": "group already exists"})
		return
	} else {
		log.Println("The group doesn't exist before")
	}

	registrationAns, err := Registryclient.Register(context.Background(), &proto.Rinfo{
		MulticastId:   id,
		MulticastType: multicastType,
		ClientPort:    uint32(GrpcPort),
	})
	if err != nil {
		g.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	members := make(map[string]Member, 0)

	for memberId, member := range registrationAns.GroupInfo.Members {
		members[memberId] = Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &MulticastGroup{
		ClientId: registrationAns.ClientId,
		Group: &MulticastInfo{
			MulticastId:      registrationAns.GroupInfo.MulticastId,
			MulticastType:    proto.MulticastType_name[int32(registrationAns.GroupInfo.MulticastType)],
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		Messages:  make([]Message, 0),
		MessageMu: sync.RWMutex{},
	}

	MulticastGroups[registrationAns.GroupInfo.MulticastId] = group

	go InitGroup(registrationAns.GroupInfo, group, len(registrationAns.GroupInfo.Members) == 1)

	g.Status(http.StatusOK)
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

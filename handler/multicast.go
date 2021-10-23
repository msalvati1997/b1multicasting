package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	context "golang.org/x/net/context"
	"net/http"
	"sync"
	"time"
)

var (
	RegClient       protoregistry.RegistryClient
	GMu             sync.RWMutex
	Delay           uint
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
	MulticastType protoregistry.MulticastType
}

type GroupConfig struct {
	MulticastType string `json:"multicast_type"`
}

// GetGroups godoc
// @Summary Get Multicast Group
// @Description Get Multicast Group
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 201 {object} MulticastGroups
//     Responses:
//       201: body:PositionResponseBody
// @Router /groups [get]
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
// CreateGroup initializes a new multicast group.
func CreateGroup(ctx *gin.Context) {
	var config GroupConfig

	multicastId := ctx.Param("multicastId")

	err := ctx.BindJSON(&config)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	multicastType, ok := registryservice.MulticastType[config.MulticastType]

	if !ok {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "multicast_type not supported", "supported": registryservice.MulticastTypes})
		return
	}

	GMu.Lock()
	defer GMu.Unlock()

	group, ok := MulticastGroups[multicastId]

	if ok {
		ctx.IndentedJSON(http.StatusConflict, gin.H{"error": "group already exists"})
		return
	}

	registrationAns, err := RegClient.Register(context.Background(), &protoregistry.Rinfo{
		MulticastId:   multicastId,
		MulticastType: multicastType,
		ClientPort:    uint32(GrpcPort),
	})

	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			MulticastType:    protoregistry.MulticastType_name[int32(registrationAns.GroupInfo.MulticastType)],
			ReceivedMessages: 0,
			Status:           protoregistry.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]Message, 0),
		GroupMu:  sync.RWMutex{},
	}

	MulticastGroups[registrationAns.GroupInfo.MulticastId] = group

	go InitGroup(registrationAns.GroupInfo, group, len(registrationAns.GroupInfo.Members) == 1)

	ctx.Status(http.StatusOK)
}

func InitGroup(info *protoregistry.MGroup, group *MulticastGroup, b bool) {
	// Waiting that the group is ready
	update(info, group)
	groupInfo, _ := StatusChange(info, group, protoregistry.Status_OPENING)

	// Communicating to the registry that the node is ready
	groupInfo, _ = RegClient.Ready(context.Background(), &protoregistry.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.ClientId,
	})

	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, protoregistry.Status_STARTING)

}

func StatusChange(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup, status protoregistry.Status) (*protoregistry.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = RegClient.GetStatus(context.Background(), &protoregistry.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}

		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == protoregistry.Status_CLOSED || groupInfo.Status == protoregistry.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

func update(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup) {
	multicastGroup.GroupMu.Lock()
	defer multicastGroup.GroupMu.Unlock()

	multicastGroup.Group.Status = protoregistry.Status_name[int32(groupInfo.Status)]

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

package multicastapp

import (
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	context "golang.org/x/net/context"
	"net/http"
	"sync"
)

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
		group.groupMu.RLock()
		groups = append(groups, group.group)
		group.groupMu.RUnlock()
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

	registrationAns, err := registryClient.Register(context.Background(), &protoregistry.Rinfo{
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
		clientId: registrationAns.ClientId,
		group: &MulticastInfo{
			MulticastId:      registrationAns.GroupInfo.MulticastId,
			MulticastType:    protoregistry.MulticastType_name[int32(registrationAns.GroupInfo.MulticastType)],
			ReceivedMessages: 0,
			Status:           protoregistry.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		messages: make([]Message, 0),
		groupMu:  sync.RWMutex{},
	}

	MulticastGroups[registrationAns.GroupInfo.MulticastId] = group

	go InitGroup(registrationAns.GroupInfo, group, len(registrationAns.GroupInfo.Members) == 1)

	ctx.Status(http.StatusOK)
}

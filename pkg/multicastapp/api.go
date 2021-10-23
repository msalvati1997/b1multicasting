package multicastapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	context "golang.org/x/net/context"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	timeout = time.Second
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
	context_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var config GroupConfig

	multicastId := ctx.Param("multicastId")

	err := ctx.BindJSON(&config)

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	multicastType, ok := registryservice.MulticastType[config.MulticastType]

	if !ok {
		response(ctx, ok, errors.New("Multicast type not supported"))
	}

	GMu.Lock()
	defer GMu.Unlock()

	group, ok := MulticastGroups[multicastId]

	if ok {
		response(ctx, ok, errors.New("Group already exists"))
	}

	registrationAns, err := registryClient.Register(context_, &protoregistry.Rinfo{
		MulticastId:   multicastId,
		MulticastType: multicastType,
		ClientPort:    uint32(GrpcPort),
	})

	if err != nil {
		response(ctx, err.Error(), errors.New("Error in registering client "))
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

	response(ctx, MulticastGroups, nil)
}
func response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		log.Println("Server Error Occured:", err)
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"data": data, "error": errorMessage})
}

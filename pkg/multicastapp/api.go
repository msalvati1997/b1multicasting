package multicastapp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	context "golang.org/x/net/context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout = time.Second
)

// @Paths Information

// GetGroups godoc
// @Summary Get Multicast Groups
// @Description Get Multicast Groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {object} MulticastInfo
// @Failure 500 {object} Response
// @Router /groups [get]
func GetGroups(g *gin.Context) {

	groups := make([]*MulticastInfo, 0)

	GMu.RLock()
	defer GMu.RUnlock()

	for _, group := range MulticastGroups {
		group.groupMu.RLock()
		groups = append(groups, group.Group)
		group.groupMu.RUnlock()
	}

	response(g, groups, nil)
}

// CreateOrJoinGroup godoc
// @Summary Create Multicast Group or join in an existing group
// @Description Create Multicast Group
// @Tags groups
// @Accept  json
// @Produce  json
// @Param request body MulticastReq true "Specify the id and type of new multicast group"
// @Success 201 {object} MulticastInfo
// @Failure 500 {object} Response
// @Router /groups [post]
func CreateOrJoinGroup(ctx *gin.Context) {
	context_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var req MulticastReq
	err := ctx.BindJSON(&req)

	multicastId := req.MulticastId

	if err != nil {
		response(ctx, nil, err)
	}
	multicastType, ok := registryservice.MulticastType[req.MulticastType]
	if !ok {
		response(ctx, ok, errors.New("Multicast type not supported"))
	}
	log.Println("Creating/Joining at multicast with type ", multicastType)

	GMu.Lock()
	defer GMu.Unlock()

	group, ok := MulticastGroups[multicastId]

	if ok {
		response(ctx, ok, errors.New("Group already exists"))
	}

	registrationAns, err := RegClient.Register(context_, &protoregistry.Rinfo{
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
		Group: &MulticastInfo{
			MulticastId:      registrationAns.GroupInfo.MulticastId,
			MulticastType:    multicastType.String(),
			ReceivedMessages: 0,
			Status:           protoregistry.Status_name[int32(registrationAns.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]Message, 0),
		groupMu:  sync.RWMutex{},
	}

	MulticastGroups[registrationAns.GroupInfo.MulticastId] = group

	go InitGroup(registrationAns.GroupInfo, group)

	response(ctx, group.Group, nil)
}

// GetGroupById godoc
// @Summary Get Multicast GroupInfo by id
// @Description Get Multicast GroupInfo by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Success 200 {object} MulticastInfo
// @Failure 500 {object} Response
// @Router /groups/{mId} [get]
// GetGroupById retrives group info by an id.
func GetGroupById(ctx *gin.Context) {

	mId := ctx.Param("mId")
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	response(ctx, group.Group, nil)
}

// DeleteGroup godoc
// @Summary Delete an existing group
// @Description  Delete an existing group
// @Tags groups
// @Accept  json
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Success 200 {object} MulticastInfo
// @Failure 500 {object} Response
// @Router /groups/{mId} [delete]
// GetGroupById  Delete an existing group
func DeleteGroup(ctx *gin.Context) {

	mId := ctx.Param("mId")
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		response(ctx, gin.H{"error": "group not found"}, errors.New("Groups not found"))
	}
	if protoregistry.Status(protoregistry.Status_value[group.Group.Status]) != protoregistry.Status_ACTIVE {
		response(ctx, gin.H{"error": "group not active"}, errors.New("The group isn't active"))
	}

	msg := basic.NewMessage(make(map[string]string), []byte("CloseGroup"))
	msg.MessageHeader["type"] = "CloseGroup"
	msg.MessageHeader["MulticastId"] = group.Group.MulticastId
	msg.MessageHeader["ClientId"] = group.clientId
	msg.MessageHeader["ProcessId"] = strconv.Itoa(utils.Myid)
	err := multicasting.Cnn.BMulticast(group.Group.MulticastId, msg)

	if err != nil {
		response(ctx, gin.H{"error": "closing group  request  error"}, err)
	}
}

// StartGroup godoc
// @Summary Start multicast group by id
// @Description Start multicast group by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Success 200 {object} MulticastInfo
// @Failure 500 {object} Response
// @Router /groups/{mId} [put]
func StartGroup(ctx *gin.Context) {

	mId := ctx.Param("mId")

	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		response(ctx, ok, errors.New("The groups doesn't exist"))
	}

	info, err := RegClient.StartGroup(context.Background(), &protoregistry.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.clientId})

	if err != nil {
		log.Println("Error in start group ", err.Error())
		return
	}

	response(ctx, info, nil)
}

// MulticastMessage godoc
// @Summary Multicast a message to a group mId
// @Description Multicast a message to a group mId
// @Tags messaging
// @Accept  json
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Param message body Message true "Message to multicast"
// @Success 200 {object} Message
// @Failure 500 {object} Response
// @Router /messaging/{mId} [POST]
func MulticastMessage(ctx *gin.Context) {
	mId := ctx.Param("mId")
	var req Message
	err := ctx.BindJSON(&req)
	if err != nil {
		response(ctx, "Error in input ", err)
	}

	group, ok := MulticastGroups[mId]
	if !ok {
		response(ctx, ok, errors.New("The groups "+mId+" doesn't exist"))
	}

	group.groupMu.RLock()
	defer group.groupMu.RUnlock()

	if protoregistry.Status(protoregistry.Status_value[group.Group.Status]) != protoregistry.Status_ACTIVE {
		ctx.IndentedJSON(http.StatusTooEarly, gin.H{"error": "group not ready"})
		return
	}
	log.Println("Trying to multicasting message to group ", mId)
	multicastType := group.Group.MulticastType
	payload := req.Payload
	mtype, ok := registryservice.MulticastType[multicastType]
	if !ok {
		response(ctx, ok, errors.New("Multicast type not supported"))
	}
	log.Println("Trying to sending ", payload)

	msg := basic.NewMessage(make(map[string]string), payload)
	msg.MessageHeader["Tranport"] = "HTTP"
	if mtype.Number() == 0 {
		msg.MessageHeader["type"] = "B"
		msg.MessageHeader["GroupId"] = group.Group.MulticastId
	}
	if mtype.Number() == 1 {
		msg.MessageHeader["type"] = "TOC"
		msg.MessageHeader["GroupId"] = group.Group.MulticastId
	}
	if mtype.Number() == 2 {
		msg.MessageHeader["type"] = "TOD"
		msg.MessageHeader["GroupId"] = group.Group.MulticastId
	}
	if mtype.Number() == 3 {
		msg.MessageHeader["type"] = "CO"
		msg.MessageHeader["GroupId"] = group.Group.MulticastId
	}

	utils2.GoPool.MessageCh <- msg

	var m Message
	m.MessageHeader = msg.MessageHeader
	m.Payload = msg.Payload

	response(ctx, m, nil)
}

// RetrieveMessages godoc
// @Summary Get Messages of a Group
// @Description Get Messages of a Group
// @Tags messaging
// @Accept  json
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Success 200 {object} []Message
// @Failure 500 {object} Response
// @Router /messaging/{mId} [get]
// GetGroupById retrieve group msg by an id
func RetrieveMessages(ctx *gin.Context) {
	mId := ctx.Param("mId")
	group, ok := MulticastGroups[mId]
	if !ok {
		response(ctx, ok, errors.New("The groups "+mId+" doesn't exist"))
	}
	response(ctx, group.Messages, nil)
}

// RetrieveDeliverQueue godoc
// @Summary Get Deliver-Message queue
// @Description Get Deliver-Message queue of Group by id
// @Tags deliver
// @Produce  json
// @Success 200 {object} []utils.Delivery
// @Failure 500 {object} Response
// @Router /deliver/ [get]
// RetrieveDeliverQueue retrieve deliver message queue
func RetrieveDeliverQueue(c *gin.Context) {

	response(c, utils2.Del.DelivererNodes, nil)
}

// RetrieveDeliverQueueOfAGroup of a specific group godoc
// @Summary Get Deliver-Message queue
// @Description Get Deliver-Message queue of Group by id
// @Tags deliver
// @Produce  json
// @Param mId path string true "Multicast group id group"
// @Success 200 {object} []utils.Delivery
// @Failure 500 {object} Response
// @Router /deliver/{mId} [get]
// RetrieveDeliverQueueOfAGroup retrieve deliver message queue of a specific group
func RetrieveDeliverQueueOfAGroup(ctx *gin.Context) {

	mId := ctx.Param("mId")
	_, ok := MulticastGroups[mId]
	if !ok {
		response(ctx, ok, errors.New("The groups "+mId+" doesn't exist"))
	}
	var groupDeliver []utils2.Delivery
	for i := 0; i < len(utils2.Del.DelivererNodes); i++ {
		if utils2.Del.DelivererNodes[i].M.MessageHeader["GroupId"] == mId {
			groupDeliver = append(groupDeliver, *utils2.Del.DelivererNodes[i])
		}
	}
	if len(groupDeliver) == 0 {
		response(ctx, gin.H{"GroupInfo": "There aren't delivered message from this group"}, nil)
	}
	response(ctx, groupDeliver, nil)
}

func response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
		c.IndentedJSON(statusCode, gin.H{"data": data, "error": errorMessage})
	} else {
		c.IndentedJSON(statusCode, gin.H{"data": data})
	}
}

type Response struct {
	H map[string]interface{}
}

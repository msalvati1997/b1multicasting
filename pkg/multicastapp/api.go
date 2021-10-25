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
	"strings"
	"sync"
	"time"
)

var (
	timeout = time.Second
)

// @BasePath /multicast/v1

// GetGroups godoc
// @Summary Get Multicast Group
// @Description Get Multicast Group
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 201 {object} MulticastInfo
//     Responses:
//       201: body:PositionResponseBody
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

// @BasePath /multicast/v1

// CreateGroup godoc
// @Summary Create Multicast Group
// @Description Create Multicast Group
// @Tags groups
// @Accept  json
// @Produce  json
// @Params  MulticastReq
//       description :
//            id of multicast and type
// @Success 201 {object} MulticastInfo
// @Router /groups [post]
// CreateGroup initializes a new multicast group or join in an group.
func CreateGroup(ctx *gin.Context) {
	context_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var req MulticastReq
	err := ctx.BindJSON(&req)

	multicastId := req.MulticastId

	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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

	go InitGroup(registrationAns.GroupInfo, group, len(registrationAns.GroupInfo.Members) == 1)

	response(ctx, group.Group, nil)
}

// @BasePath /multicast/v1

// GetGroupById godoc
// @Summary Get Multicast Group by id
// @Description Get Multicast Group by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Params MulticastId
//       Description : id of multicast
// @Success 201 {object} MulticastInfo
// @Router /groups/:mId [get]
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

	group.groupMu.RLock()
	defer group.groupMu.RUnlock()

	response(ctx, group.Group, nil)
}

// @BasePath /multicast/v1

// StartGroup godoc
// @Summary Start multicast by id
// @Description Start multicast by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 201 {object} MulticastInfo
// @Router /groups/:mId [put]
// StartGroup starting multicast group
func StartGroup(ctx *gin.Context) {

	mId := ctx.Param("mId")

	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[mId]

	if !ok {
		response(ctx, ok, errors.New("The groups doesn't exist"))
	}

	groupInfo, err := registryClient.StartGroup(context.Background(), &protoregistry.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.clientId})

	if err != nil {
		log.Println("Error in start group ", err.Error())
		return
	}

	log.Println("Group ", groupInfo.MulticastId, "start with types of communication ", groupInfo.MulticastType)

	if groupInfo.MulticastType.String() == "BMULTICAST" {
		log.Println("STARTING BMULTICAST COMMUNICATION")
	}
	if groupInfo.MulticastType.String() == "TOCMULTICAST" {
		log.Println("STARTING TOC COMMUNICATION")
		members := []string{}
		for k := range groupInfo.Members {
			members = append(members, k)
		}
		sequencerPort := multicasting.SelectingSequencer(members)
		seqCon, err := multicasting.Cnn.GetGrpcClient(sequencerPort)
		if err != nil {
			log.Println("Error in find connection with sequencer..", err.Error())
			return
		}
		seq := false
		if utils.MyAdress == sequencerPort {
			log.Println("I'm the sequencer of MulticastGroup", groupInfo.MulticastId)
			seq = true
		} else {
			log.Println("The sequencer nodes is at port", sequencerPort)
		}
		multicasting.Seq.MulticastId = groupInfo.MulticastId
		multicasting.Seq.SeqConn = *seqCon
		multicasting.Seq.Conns = multicasting.Cnn
		multicasting.Seq.B = seq
		multicasting.Seq.SeqPort = sequencerPort
		go utils2.TOCDeliver()
	}
	if groupInfo.MulticastType.String() == "TODMULTICAST" {
		log.Println("STARTING TOD COMMUNICATION")
		go utils2.TODDeliver()
	}
	if groupInfo.MulticastType.String() == "COMULTICAST" {
		utils.Vectorclock = utils.NewVectorClock(len(groupInfo.Members))
		log.Println("STARTING CO COMMUNICATION")
		go utils2.CODeliver()
	}

	response(ctx, group.Group, nil)
}

// @BasePath /multicast/v1

// MulticastMessage godoc
// @Summary Multicast a message to a group G
// @Description Multicast a message to a group G
// @Tags messaging
// @Accept  json
// @Produce  json
// @Success 201 {object} Message
// @Router /messaging/:mId [POST]
// MulticastMessage Multicast a message to a group mId
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
// @Summary Get Message of Group by id
// @Description Get Message of Group by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Params Message
//       Description : id of multicast
// @Success 201 {object} []Message
// @Router /messaging/:mId [get]
// GetGroupById retrieve group msg by an id
func RetrieveMessages(ctx *gin.Context) {
	mId := ctx.Param("mId")
	group, ok := MulticastGroups[mId]
	if !ok {
		response(ctx, ok, errors.New("The groups "+mId+" doesn't exist"))
	}
	response(ctx, group.Messages, nil)
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

// RetrieveDeliverQueue godoc
// @Summary Get Deliver-Message queue
// @Description Get Deliver-Message of Group by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Params Message
//       Description : id of multicast
// @Success 201 {object} []Message
// @Router /deliver/:mId [get]
// RetrieveDeliverQueue retrieve deliver message queue
func RetrieveDeliverQueue(c *gin.Context) {
	response(c, utils2.Del.DelivererNodes, nil)
}

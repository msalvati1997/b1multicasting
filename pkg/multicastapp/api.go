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
	"google.golang.org/grpc/peer"
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

// @BasePath /multicast/api

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
// @Success 201 {object} MulticastInfo
//     Responses:
//       201: body:PositionResponseBody
// @Router /groups [post]
// CreateGroup initializes a new multicast group or join in an group.
func CreateGroup(ctx *gin.Context) {
	context_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	source, ok := peer.FromContext(ctx)
	src := source.Addr.String()
	srcAddr := src[:strings.LastIndexByte(src, ':')]
	myid := strings.Split(srcAddr, ":")
	mid := myid[len(myid)-1]
	i, _ := strconv.Atoi(mid)
	utils.Myid = i
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

	response(ctx, group.group, nil)
}

// GetGroupById godoc
// @Summary Get Multicast Group by id
// @Description Get Multicast Group by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Param post body MulticastId
// @Success 201 {object} MulticastInfo
//     Responses:
//       201: body:PositionResponseBody
// @Router /groups/:mId [get]
// GetGroupById retrives group info by an id.
func GetGroupById(ctx *gin.Context) {
	multicastId := ctx.Param("multicastId")

	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[multicastId]

	if !ok {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	group.groupMu.RLock()
	defer group.groupMu.RUnlock()

	response(ctx, group.group, nil)
}

// StartGroup godoc
// @Summary Start multicast by id
// @Description Start multicast by id
// @Tags groups
// @Accept  json
// @Produce  json
// @Param post body MulticastId
// @Success 201 {object} MulticastInfo
//     Responses:
//       201: body:PositionResponseBody
// @Router /groups/:mId [put]
// StartGroup starting multicast group
func StartGroup(ctx *gin.Context) {
	multicastId := ctx.Param("multicastId")

	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[multicastId]

	if !ok {
		response(ctx, ok, errors.New("The groups doesn't exist"))
	}
	source, ok := peer.FromContext(ctx)
	src := source.Addr.String()
	srcAddr := src[:strings.LastIndexByte(src, ':')]
	myid := strings.Split(srcAddr, ":")
	mid := myid[len(myid)-1]
	i, _ := strconv.Atoi(mid)
	utils.Myid = i
	groupInfo, err := registryClient.StartGroup(context.Background(), &protoregistry.RequestData{
		MulticastId: group.group.MulticastId,
		MId:         group.clientId})

	if err != nil {
		log.Println("Error in start group ", err.Error())
		return
	}
	log.Println("Group starting..")

	if groupInfo.MulticastType.Number() == 0 {
		log.Println("STARTING BMULTICAST COMMUNICATION")
	}
	if groupInfo.MulticastType.Number() == 1 {
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
		if srcAddr == sequencerPort {
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
	if groupInfo.MulticastType.Number() == 2 {
		log.Println("STARTING TOD COMMUNICATION")
		go utils2.TODDeliver()
	}
	if groupInfo.MulticastType.Number() == 3 {
		utils.Vectorclock = utils.NewVectorClock(len(groupInfo.Members))
		log.Println("STARTING CO COMMUNICATION")
		go utils2.CODeliver()
	}

	response(ctx, group.group, nil)
}

// MulticastMessage godoc
// @Summary Multicast a message to a group G
// @Description Multicast a message to a group G
// @Tags messaging
// @Accept  json
// @Produce  json
// @Param post body MessageRequest
// @Success 201 {object}
//     Responses:  basic.Message
//       201: body:PositionResponseBody
// @Router /messaging/:mId [POST]
// MulticastMessage Multicast a message to a group G
func MulticastMessage(ctx *gin.Context) {
	var req MessageRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		return
	}
	group, ok := MulticastGroups[req.multicastId.MulticastId]
	if !ok {
		response(ctx, ok, errors.New("The groups doesn't exist"))
	}
	multicastType := group.group.MulticastType
	payload := req.message.Payload
	mtype, ok := registryservice.MulticastType[multicastType]
	if !ok {
		response(ctx, ok, errors.New("Multicast type not supported"))
	}
	msg := basic.NewMessage(make(map[string]string), payload)
	if mtype.Number() == 0 {
		msg.MessageHeader["type"] = "B"
		msg.MessageHeader["GroupId"] = group.group.MulticastId
	}
	if mtype.Number() == 1 {
		msg.MessageHeader["type"] = "TOC"
		msg.MessageHeader["GroupId"] = group.group.MulticastId
	}
	if mtype.Number() == 2 {
		msg.MessageHeader["type"] = "TOD"
		msg.MessageHeader["GroupId"] = group.group.MulticastId
	}
	if mtype.Number() == 3 {
		msg.MessageHeader["type"] = "CO"
		msg.MessageHeader["GroupId"] = group.group.MulticastId
	}
	utils2.GoPool.MessageCh <- msg
	response(ctx, msg, nil)
}

func response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"data": data, "error": errorMessage})
}

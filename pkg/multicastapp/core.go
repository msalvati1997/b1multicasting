package multicastapp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/docs"
	utils2 "github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

var (
	RegClient       protoregistry.RegistryClient
	GMu             sync.RWMutex
	Delay           uint
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
	Application     bool
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	clientId  string
	Group     *MulticastInfo
	groupMu   sync.RWMutex
	Messages  []Message
	MessageMu sync.RWMutex
}

type Message struct {
	MessageHeader map[string]string `json:"MessageHeader"`
	Payload       []byte            `json:"Payload"`
}

type MessageRequest struct {
	MulticastId string `json:"multicast_id"`
	message     Message
}

type MulticastInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

type MulticastReq struct {
	MulticastId   string `json:"multicast_id"`
	MulticastType string `json:"multicast_type"`
}
type MulticastId struct {
	MulticastId string `json:"multicast_id"`
}

type GroupConfig struct {
	MulticastType string `json:"multicast_type"`
}
type routes struct {
	router *gin.Engine
}

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	GrpcPort = grpcP
	Delay = dl

	var err error
	Application = true
	utils2.Vectorclock = utils2.NewVectorClock(4)

	RegClient, err = client.Connect(registryAddr)

	if err != nil {
		return err
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	utils.GoPool.Initialize(int(numThreads))

	r := routes{
		router: gin.Default(),
	}
	docs.SwaggerInfo.BasePath = relativePath
	docs.SwaggerInfo.Schemes = []string{"http"}
	v1 := r.router.Group(relativePath)
	r.addGroups(v1)
	r.addMessaging(v1)
	r.addSwagger(v1)
	r.addDeliver(v1)
	err = r.router.Run(fmt.Sprintf(":%d", restPort))
	return err
}

func (r routes) addGroups(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	groups.GET("/", GetGroups)
	groups.POST("/", CreateOrJoinGroup)
	groups.PUT("/:mId", StartGroup)
	groups.GET("/:mId", GetGroupById)
	groups.DELETE("/:mId", DeleteGroup)
}

func (r routes) addMessaging(rg *gin.RouterGroup) {
	groups := rg.Group("/messaging")
	groups.POST("/:mId", MulticastMessage)
	groups.GET("/:mId", RetrieveMessages)
}

func (r routes) addSwagger(rg *gin.RouterGroup) {
	groups := rg.Group("/swagger")
	groups.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r routes) addDeliver(rg *gin.RouterGroup) {
	groups := rg.Group("/deliver")
	groups.GET("/:mId", RetrieveDeliverQueueOfAGroup)
	groups.GET("/", RetrieveDeliverQueue)
}

func InitGroup(info *protoregistry.MGroup, group *MulticastGroup) {
	// Waiting that the group is ready
	log.Println("Waiting for the group to be ready")

	update(info, group)
	groupInfo, err := StatusChange(info, group, protoregistry.Status_OPENING)
	if err != nil {
		return
	}

	log.Println("Group ready, initializing multicast")
	// Initializing  data structures
	err = initialiazeGroupCommunication(info, group)

	if err != nil {
		return
	}

	groupInfo, err = RegClient.Ready(context.Background(), &protoregistry.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.clientId,
	})
	if err != nil {
		return
	}

	log.Println("Waiting for the other nodes")
	if groupInfo.MulticastType.String() == "COMULTICAST" && utils.G3 == false {
		go utils.CODeliver()
		utils.G3 = true
	}
	if groupInfo.MulticastType.String() == "TOCMULTICAST" && utils.G2 == false {
		go utils.TOCDeliver()
		utils.G2 = true
	}
	if groupInfo.MulticastType.String() == "TODMULTICAST" && utils.G1 == false {
		go utils.TODDeliver()
		utils.G1 = true
	}
	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, protoregistry.Status_STARTING)

	if err != nil {
		return
	}

	log.Println("Ready to multicast")

}

func initialiazeGroupCommunication(groupInfo *protoregistry.MGroup, group *MulticastGroup) error {

	var members []string

	for memberId, member := range group.Group.Members {
		if memberId != group.clientId {
			members = append(members, member.Address)
		}
	}
	members = append(members, group.clientId)
	_, err := multicasting.Connections(members, int(Delay))
	if err != nil {
		return err
	}

	if groupInfo.MulticastType.String() == "BMULTICAST" {
		log.Println("STARTING BMULTICAST COMMUNICATION")
	}
	if groupInfo.MulticastType.String() == "TOCMULTICAST" {
		log.Println("STARTING TOC COMMUNICATION")
		multicasting.Seq.MulticastId = groupInfo.MulticastId
		multicasting.Seq.Conns = multicasting.Cnn
		sequencerPort := multicasting.SelectingSequencer(members, true)
		log.Println("Sequencer selected with port ", sequencerPort)
		multicasting.Seq.SeqPort = sequencerPort
		seqconnection, err := multicasting.Cnn.GetGrpcClient(sequencerPort)
		if err != nil {
			log.Println("Error in find connection with sequencer..", err.Error())
		}
		multicasting.Seq.SeqConn = *seqconnection
		log.Println("The sequencer nodes is at port", seqconnection.GetTarget())
	}

	if groupInfo.MulticastType.String() == "TODMULTICAST" {
		log.Println("STARTING TOD COMMUNICATION")
	}
	if groupInfo.MulticastType.String() == "COMULTICAST" {
		log.Println("STARTING CO COMMUNICATION")
	}
	return nil
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
	multicastGroup.groupMu.Lock()
	defer multicastGroup.groupMu.Unlock()

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

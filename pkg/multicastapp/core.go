package multicastapp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"golang.org/x/net/context"
	"log"
	"sync"
	"time"
)

var (
	registryClient  protoregistry.RegistryClient
	GMu             sync.RWMutex
	Delay           uint
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	clientId  string
	group     *MulticastInfo
	groupMu   sync.RWMutex
	messages  []Message
	messageMu sync.RWMutex
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
	Error            error             `json:"error"`
}

// ErrorResponse defines an error response returned upon any request failure.
type ErrorResponse struct {
	Error string `json:"error"`
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

	registryClient, err = client.Connect(registryAddr)

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
	v1 := r.router.Group(relativePath)
	r.addGroups(v1)

	err = r.router.Run(fmt.Sprintf(":%d", restPort))
	return err
}

func (r routes) addGroups(rg *gin.RouterGroup) {
	groups := rg.Group("/groups")
	groups.GET("/", GetGroups)
	groups.POST("/", CreateGroup)
	groups.PUT("/:mId", StartGroup)
	groups.GET("/:mId", GetGroupById)
}

func InitGroup(info *protoregistry.MGroup, group *MulticastGroup, b bool) {
	// Waiting that the group is ready
	log.Println("Waiting for the group to be ready")

	update(info, group)
	groupInfo, err := StatusChange(info, group, protoregistry.Status_OPENING)
	if err != nil {
		group.groupMu.Lock()
		group.group.Error = err
		group.groupMu.Unlock()
		return
	}

	log.Println("Group ready, initializing multicast")

	// Initializing  data structures
	err = initializeMulticast(group, b)

	if err != nil {
		group.groupMu.Lock()
		group.group.Error = err
		group.groupMu.Unlock()
		return
	}

	log.Println("Notify the registry that the multicaster is ready")
	// Communicating to the registry that the node is ready
	groupInfo, err = registryClient.Ready(context.Background(), &protoregistry.RequestData{
		MulticastId: group.group.MulticastId,
		MId:         group.clientId,
	})
	if err != nil {
		group.groupMu.Lock()
		group.group.Error = err
		group.groupMu.Unlock()
		return
	}

	log.Println("Waiting for the other nodes")
	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, protoregistry.Status_STARTING)

	if err != nil {
		group.groupMu.Lock()
		group.group.Error = err
		group.groupMu.Unlock()
		return
	}

	log.Println("Ready to multicast")
}

func initializeMulticast(group *MulticastGroup, b bool) error {

	var members []string
	//effettuo la connessione degli altri nodi come clients
	for memberId, member := range group.group.Members {
		if memberId != group.clientId {
			log.Println("Connecting to: %s", member.Address)
			members = append(members, member.Address)
		}
	}
	members = append(members, group.clientId)
	_, err := multicasting.Connections(members, int(Delay))
	if err != nil {
		return err
	}

	return nil
}

func StatusChange(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup, status protoregistry.Status) (*protoregistry.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = registryClient.GetStatus(context.Background(), &protoregistry.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}
		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == protoregistry.Status_CLOSED || groupInfo.Status == protoregistry.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}
	log.Println("Group status changes")

	return groupInfo, nil
}

func update(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup) {
	multicastGroup.groupMu.Lock()
	defer multicastGroup.groupMu.Unlock()

	multicastGroup.group.Status = protoregistry.Status_name[int32(groupInfo.Status)]

	for clientId, member := range groupInfo.Members {
		m, ok := multicastGroup.group.Members[clientId]

		if !ok {
			m = Member{
				MemberId: member.Id,
				Address:  member.Address,
				Ready:    member.Ready,
			}

			multicastGroup.group.Members[clientId] = m
		}

		if m.Ready != member.Ready {
			m.Ready = member.Ready
			multicastGroup.group.Members[clientId] = m
		}

	}
}

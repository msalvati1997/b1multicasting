package multicastapp

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"golang.org/x/net/context"
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

type GroupConfig struct {
	MulticastType string `json:"multicast_type"`
}

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	GrpcPort = grpcP
	Delay = dl

	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)

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

	router := gin.Default()

	routerGroup := router.Group(relativePath)

	routerMap, err := GetRouterMethods(routerGroup)

	for path, methods := range GetRestAPI() {
		for methodType, method := range methods {
			routerMethod, ok := routerMap[methodType]

			if !ok {
				return errors.New(fmt.Sprintf("missing routing method: %s", HttpMethodName[methodType]))
			}

			routerMethod(path, method)

		}
	}
	err = router.Run(fmt.Sprintf(":%d", restPort))
	if err != nil {
		return err
	}
	return err
}

var restAPI = map[string]map[HttpMethod]func(ctx *gin.Context){

	"/groups": {
		HttpMethodGET: GetGroups,
	},
	"/groups/:multicastId": {
		HttpMethodPUT: CreateGroup,
	},
}

func GetRestAPI() map[string]map[HttpMethod]func(ctx *gin.Context) {
	api := make(map[string]map[HttpMethod]func(ctx *gin.Context))

	for key, value := range restAPI {
		api[key] = value
	}

	return api
}

const (
	HttpMethodGET HttpMethod = iota
	HttpMethodPOST
	HttpMethodPUT
	HttpMethodDELETE
	HttpMethodHEAD
)

var (
	HttpMethodName = map[HttpMethod]string{
		HttpMethodGET:    "GET",
		HttpMethodPOST:   "POST",
		HttpMethodPUT:    "PUT",
		HttpMethodDELETE: "DELETE",
		HttpMethodHEAD:   "HEAD",
	}
	MulticastTypeValue = map[string]HttpMethod{
		"GET":    HttpMethodGET,
		"POST":   HttpMethodPOST,
		"PUT":    HttpMethodPUT,
		"DELETE": HttpMethodDELETE,
		"HEAD":   HttpMethodHEAD,
	}
)

func GetRouterMethods(router *gin.RouterGroup) (map[HttpMethod]func(path string, handlers ...gin.HandlerFunc) gin.IRoutes, error) {
	routerMap := make(map[HttpMethod]func(path string, handlers ...gin.HandlerFunc) gin.IRoutes)

	routerMap[HttpMethodGET] = router.GET
	routerMap[HttpMethodPOST] = router.POST
	routerMap[HttpMethodPUT] = router.PUT
	return routerMap, nil

}

type HttpMethod int

func GroupsApi(router *gin.RouterGroup) {
	router.GET("/", GetGroups)
	router.POST("/", CreateGroup)
}

func InitGroup(info *protoregistry.MGroup, group *MulticastGroup, b bool) {
	// Waiting that the group is ready
	update(info, group)
	groupInfo, _ := StatusChange(info, group, protoregistry.Status_OPENING)

	// Communicating to the registry that the node is ready
	groupInfo, _ = registryClient.Ready(context.Background(), &protoregistry.RequestData{
		MulticastId: group.group.MulticastId,
		MId:         group.clientId,
	})

	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, protoregistry.Status_STARTING)

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

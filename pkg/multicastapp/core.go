package multicastapp

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/context"
	"log"
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

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool, M cmux.CMux) error {
	GrpcPort = grpcP
	Delay = dl
	httpL := M.Match(cmux.HTTP1Fast())
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		RegClient, err = client.Connect(registryAddr)
		if err != nil {
			log.Println("Error in connecting registry client ", err.Error())
		}
		wg.Done()
	}()
	wg.Wait()

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	utils.GoPool.Initialize(int(numThreads))

	router := gin.Default()

	routerGroup := router.Group(relativePath)
	routerGroup.GET("/groups", GetGroups)
	routerGroup.POST("/groups", CreateGroup)

	go func() {
		err := router.RunListener(httpL)
		if err != nil {
			return
		}
	}()
	return err
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

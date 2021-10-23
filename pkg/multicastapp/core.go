// Package distribuitedmulticast provides functions to run a Http REST server running a multicast node
package multicastapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"sync"
	"time"
)

const appHeader = "multicast"

var (
	multicastGroups map[string]*MulticastGroup
	groupsMu        sync.RWMutex
	registryClient  protoregistry.RegistryClient
	grpcPort        uint
	delay           uint
)

func init() {
	multicastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	clientId  string
	group     *MulticastGroupInfo
	groupMu   sync.RWMutex
	messages  []basic.Message
	messageMu sync.RWMutex
}

type MulticastGroupInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
	Error            error             `json:"error"`
}

type RegistryGroupInfo struct {
	MulticastId   string            `json:"multicast_id"`
	MulticastType string            `json:"multicast_type"`
	Status        string            `json:"status"`
	Members       map[string]Member `json:"members"`
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

type GroupConfig struct {
	MulticastType string `json:"multicast_type"`
}

// AppCmd represents a command that can be used between multicast nodes
type AppCmd int

// Run initializes and executes the Rest HTTP server and the multicast node
func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	grpcPort = grpcP
	delay = dl
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

	return nil
}

// initializeGroup initializes a multicast group
func initializeGroup(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup, isMaster bool) {

	// Waiting that the group is ready
	updateMulticastGroup(groupInfo, multicastGroup)
	groupInfo, err := waitStatusChange(groupInfo, multicastGroup, protoregistry.Status_OPENING)

	if err != nil {
		multicastGroup.groupMu.Lock()
		multicastGroup.group.Error = err
		multicastGroup.groupMu.Unlock()
		return
	}

	if err != nil {
		multicastGroup.groupMu.Lock()
		multicastGroup.group.Error = err
		multicastGroup.groupMu.Unlock()
		return
	}

	// Communicating to the registry that the node is ready
	groupInfo, err = registryClient.Ready(context.Background(), &protoregistry.RequestData{MulticastId: multicastGroup.group.MulticastId, MId: multicastGroup.clientId})

	if err != nil {
		multicastGroup.groupMu.Lock()
		multicastGroup.group.Error = err
		multicastGroup.groupMu.Unlock()
		return
	}

	// Waiting tha all other nodes are ready
	updateMulticastGroup(groupInfo, multicastGroup)
	groupInfo, err = waitStatusChange(groupInfo, multicastGroup, protoregistry.Status_STARTING)

	if err != nil {
		multicastGroup.groupMu.Lock()
		multicastGroup.group.Error = err
		multicastGroup.groupMu.Unlock()
		return
	}

}

// waitStatusChange waits until the status of the group in the registry changes
func waitStatusChange(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup, status protoregistry.Status) (*protoregistry.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = registryClient.GetStatus(context.Background(), &protoregistry.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}

		updateMulticastGroup(groupInfo, multicastGroup)

	}

	if groupInfo.Status == protoregistry.Status_CLOSED || groupInfo.Status == protoregistry.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

// updateMulticastGroup updates local group infos with ones received from the registry
func updateMulticastGroup(groupInfo *protoregistry.MGroup, multicastGroup *MulticastGroup) {
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

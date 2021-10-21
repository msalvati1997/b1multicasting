package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	_ "flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	_ "github.com/msalvati1997/b1multicasting/pkg/basic"
	serverservice "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	serverregistry "github.com/msalvati1997/b1multicasting/pkg/registry/server"
	_ "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	RegistryClient  proto.RegistryClient
	GMu             sync.RWMutex
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
	groupsName      []MulticastId
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

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

type MulticastId struct {
	MulticastId string `json:"multicast_id"`
}

func main() {

	d := utils.GetEnvIntWithDefault("DELAY", 0)
	nt := utils.GetEnvIntWithDefault("NUM_THREADS", 1)
	verbose := utils.GetEnvStringWithDefault("VERBOSE", "ON")
	rg := utils.GetEnvBoolWithDefault("REGISTRY", false)
	app := utils.GetEnvBoolWithDefault("APP", false)
	gPort := utils.GetEnvIntWithDefault("GRPC_PORT", 90)
	rPort := utils.GetEnvIntWithDefault("REST_PORT", 80)

	delay := flag.Uint("DELAY", uint(d), "delay for sending operations (ms)")
	grpcPort := flag.Uint("GRPC_PORT", uint(gPort), "port number of the grpc server")
	restPort := flag.Uint("REST_PORT", uint(rPort), "port number of the grpc server")
	numThreads := flag.Uint("NUM_THREADS", uint(nt), "number of threads used to multicast messages")
	verb := flag.String("VERBOSE", verbose, "Turn verbose mode on or off.")
	registry_addr := flag.String("REGISTRY_ADDR", ":90", "service registry adress")
	reg := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *reg {
		services = append(services, serverregistry.Registration)
	}
	if *application {
		services = append(services, serverservice.RegisterService)
	}
	log.Println("start")
	var wg sync.WaitGroup
	wg.Add(2)
	go func(w *sync.WaitGroup) {
		err := serverservice.RunServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		w.Done()
	}(&wg)
	if *application {
		go func() {
			err := Run(*grpcPort, *restPort, *registry_addr, *numThreads, *delay, *verb)
			if err != nil {
				return
			}
			wg.Done()
		}()
	}
	wgChan := make(chan bool)

	go func() {
		wg.Wait()
		wgChan <- true
	}()

	select {
	case <-wgChan:
	}
}

// @title Orders API
// @version 1.0
// @description This is a sample service for managing groups multicast
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /
func Run(grpcP uint, restPort uint, registryAddr string, numThreads uint, dl uint, verbose string) error {
	var err error
	GrpcPort = grpcP

	if err != nil {
		return err
	}
	newRouter := mux.NewRouter()
	newRouter.HandleFunc("/groups", GetGroups).Methods("GET")
	newRouter.HandleFunc("/groups", CreateGroup).Methods("POST")

	//newRouter.HandleFunc("/group", CreateGroup).Methods("POST")
	//utils2.GoPool.Initialize(int(numThreads))
	newRouter.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", restPort), newRouter))
	return nil
}

// GetGroups godoc
// @Summary Get details of all groups
// @Description Get details of all groups
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {object} MulticastGroup
// @Router /MulticastGroups [get]
func GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MulticastGroups)
}

// CreateGroup godoc
// @Summary Create a group from an id
// @Description Create a group from an id
// @Tags groups
// @Accept  json
// @Produce  json
// @Success 200 {object} MulticastGroup
// @Router /MulticastGroup [post]
func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var multicastId MulticastId
	json.NewDecoder(r.Body).Decode(&multicastId)
	groupsName = append(groupsName, multicastId)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(multicastId)
	GMu.RLock()
	defer GMu.RUnlock()

	group, ok := MulticastGroups[multicastId.MulticastId]

	if ok {
		log.Println("The group already exist")
	}

	register, err := RegistryClient.Register(context.Background(), &proto.Rinfo{
		MulticastId: multicastId.MulticastId,
		ClientPort:  uint32(GrpcPort),
	})
	if err != nil {
		log.Println("Problem in regitering member to group")
	}
	members := make(map[string]Member, 0)

	for memberId, member := range register.GroupInfo.Members {
		members[memberId] = Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &MulticastGroup{
		ClientId: register.ClientId,
		Group: &MulticastInfo{
			MulticastId:      register.GroupInfo.MulticastId,
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(register.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]Message, 0),
	}

	MulticastGroups[register.GroupInfo.MulticastId] = group
	go InitGroup(register.GroupInfo, group, len(register.GroupInfo.Members) == 1)
	log.Println("Group created")

	w.WriteHeader(http.StatusCreated)
}

func InitGroup(info *proto.MGroup, group *MulticastGroup, b bool) {
	// Waiting that the group is ready
	update(info, group)
	groupInfo, _ := StatusChange(info, group, proto.Status_OPENING)

	// Communicating to the registry that the node is ready
	groupInfo, _ = RegistryClient.Ready(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.ClientId,
	})

	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, _ = StatusChange(groupInfo, group, proto.Status_STARTING)

}

func StatusChange(groupInfo *proto.MGroup, multicastGroup *MulticastGroup, status proto.Status) (*proto.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = RegistryClient.GetStatus(context.Background(), &proto.MulticastId{MulticastId: groupInfo.MulticastId})
		if err != nil {
			return nil, err
		}

		update(groupInfo, multicastGroup)
	}

	if groupInfo.Status == proto.Status_CLOSED || groupInfo.Status == proto.Status_CLOSING {
		return nil, errors.New("multicast group is closed")
	}

	return groupInfo, nil
}

func update(groupInfo *proto.MGroup, multicastGroup *MulticastGroup) {
	multicastGroup.GroupMu.Lock()
	defer multicastGroup.GroupMu.Unlock()

	multicastGroup.Group.Status = proto.Status_name[int32(groupInfo.Status)]

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

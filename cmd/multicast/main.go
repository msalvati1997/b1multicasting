package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	_ "flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/msalvati1997/b1multicasting/internal/app"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	serverservice "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	serverregistry "github.com/msalvati1997/b1multicasting/pkg/registry/server"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

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

// Run initializes and executes the Rest HTTP server
func Run(grpcP uint, restPort uint, registryAddr string, numThreads uint, dl uint, verbose string) error {
	var err error
	app.GrpcPort = grpcP

	if err != nil {
		return err
	}
	newRouter := mux.NewRouter()
	newRouter.HandleFunc("/CreateGroup", CreateGroup).Methods("POST")
	utils2.GoPool.Initialize(int(numThreads))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", restPort), newRouter))
	return nil
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("please provide valid data to create", err)
	}
	var multicastId string
	json.Unmarshal(reqbody, &multicastId)
	app.GMu.RLock()
	defer app.GMu.RUnlock()

	group, ok := app.MulticastGroups[multicastId]

	if ok {

	}

	register, err := app.RegistryClient.Register(context.Background(), &proto.Rinfo{
		MulticastId: multicastId,
		ClientPort:  uint32(app.GrpcPort),
	})
	if err != nil {

	}
	members := make(map[string]app.Member, 0)

	for memberId, member := range register.GroupInfo.Members {
		members[memberId] = app.Member{
			MemberId: member.Id,
			Address:  member.Address,
			Ready:    member.Ready,
		}
	}

	group = &app.MulticastGroup{
		ClientId: register.ClientId,
		Group: &app.MulticastInfo{
			MulticastId:      register.GroupInfo.MulticastId,
			ReceivedMessages: 0,
			Status:           proto.Status_name[int32(register.GroupInfo.Status)],
			Members:          members,
		},
		Messages: make([]basic.Message, 0),
	}

	app.MulticastGroups[register.GroupInfo.MulticastId] = group
	go InitGroup(register.GroupInfo, group, len(register.GroupInfo.Members) == 1)
	log.Println("Group created")
	w.WriteHeader(http.StatusCreated)
}

func InitGroup(info *proto.MGroup, group *app.MulticastGroup, b bool) {
	// Waiting that the group is ready
	update(info, group)
	groupInfo, err := StatusChange(info, group, proto.Status_OPENING)

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}
	// Communicating to the registry that the node is ready
	groupInfo, err = app.RegistryClient.Ready(context.Background(), &proto.RequestData{
		MulticastId: group.Group.MulticastId,
		MId:         group.ClientId,
	})

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}

	// Waiting tha all other nodes are ready
	update(groupInfo, group)
	groupInfo, err = StatusChange(groupInfo, group, proto.Status_STARTING)

	if err != nil {
		group.GroupMu.Lock()
		group.Group.Error = err
		group.GroupMu.Unlock()
		return
	}

}

func StatusChange(groupInfo *proto.MGroup, multicastGroup *app.MulticastGroup, status proto.Status) (*proto.MGroup, error) {
	var err error

	for groupInfo.Status == status {
		time.Sleep(time.Second * 5)
		groupInfo, err = app.RegistryClient.GetStatus(context.Background(), &proto.MulticastId{MulticastId: groupInfo.MulticastId})
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

func update(groupInfo *proto.MGroup, multicastGroup *app.MulticastGroup) {
	multicastGroup.GroupMu.Lock()
	defer multicastGroup.GroupMu.Unlock()

	multicastGroup.Group.Status = proto.Status_name[int32(groupInfo.Status)]

	for clientId, member := range groupInfo.Members {
		m, ok := multicastGroup.Group.Members[clientId]

		if !ok {
			m = app.Member{
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

package cmd

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	registry "github.com/msalvati1997/b1multicasting/pkg/registry/client"
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"log"
	"sync"
)

var (
	RegistryClient  proto.RegistryClient
	GMu             sync.RWMutex
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	ClientId  string
	Group     *MulticastInfo
	GroupMu   sync.RWMutex
	Messages  []basic.Message
	MessageMu sync.RWMutex
}

type MulticastInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
	Error            error             `json:"error"`
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

// Run initializes and executes the Rest HTTP server
func Run(grpcP uint, restPort uint, registryAddr string, numThreads uint, dl uint, verbose string) error {
	var err error
	GrpcPort = grpcP
	RegistryClient, err = registry.Connect(registryAddr)
	if err != nil {
		return err
	}
	beego.BConfig.AppName = "MulticastApp"
	if verbose == "ON" {
		beego.BConfig.Log = beego.LogConfig{AccessLogs: true}
		beego.BConfig.Log.AccessLogs = true
	}
	beego.BConfig.RunMode = "dev"
	beego.BConfig.Listen.HTTPPort = int(restPort)
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.BConfig.CopyRequestBody = true

	utils2.GoPool.Initialize(int(numThreads))
	//	go utils2.TODDeliver()
	//	go utils2.CODeliver()
	//	go utils2.TOCDeliver()
	log.Println("Run on port ", restPort)
	beego.Run(fmt.Sprintf(":%d", restPort))
	return nil
}

package multicastapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/handler"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"golang.org/x/net/context"
	"log"
	"sync"
)

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	handler.GrpcPort = grpcP
	handler.Delay = dl
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		handler.RegClient, err = client.Connect(registryAddr)
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
	multicastType, _ := registryservice.MulticastType["BMULTICAST"]

	registrationAns, err := handler.RegClient.Register(context.Background(), &protoregistry.Rinfo{
		MulticastId:   "BASICGROUP",
		MulticastType: multicastType,
		ClientPort:    uint32(handler.GrpcPort),
	})

	log.Println("Group info", registrationAns.GetGroupInfo())
	utils.GoPool.Initialize(int(numThreads))

	router := gin.Default()

	routerGroup := router.Group(relativePath)
	routerGroup.GET("/groups", handler.GetGroups)
	routerGroup.POST("/groups", handler.CreateGroup)
	err = router.Run(fmt.Sprintf(":%d", restPort))
	if err != nil {
		return err
	}

	return nil
}

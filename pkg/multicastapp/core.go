package multicastapp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/handler"
	"github.com/msalvati1997/b1multicasting/pkg/reg/client"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
)

func Run(grpcP, restPort uint, registryAddr, relativePath string, numThreads, dl uint, debug bool) error {
	handler.GrpcPort = grpcP
	handler.Delay = dl
	var err error

	handler.RegClient, err = client.Connect(registryAddr)

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
	routerGroup.GET("/groups", handler.GetGroups)
	routerGroup.POST("/groups", handler.CreateGroup)

	err = router.Run(fmt.Sprintf(":%d", restPort))
	if err != nil {
		return err
	}

	return nil
}

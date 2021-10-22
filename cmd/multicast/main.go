package main

import (
	_ "context"
	"flag"
	_ "flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/handler"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	_ "github.com/msalvati1997/b1multicasting/pkg/basic"
	basic "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	_ "github.com/msalvati1997/b1multicasting/pkg/reg"
	clientregistry "github.com/msalvati1997/b1multicasting/pkg/reg/client"
	_ "github.com/msalvati1997/b1multicasting/pkg/reg/proto"
	rgstr "github.com/msalvati1997/b1multicasting/pkg/reg/server"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

func main() {

	//	d := utils.GetEnvIntWithDefault("DELAY", 0)
	//	nt := utils.GetEnvIntWithDefault("NUM_THREADS", 1)
	//verbose := utils.GetEnvStringWithDefault("VERBOSE", "ON")
	rg := utils.GetEnvBoolWithDefault("REGISTRY", false)
	app := utils.GetEnvBoolWithDefault("APP", false)
	gPort := utils.GetEnvIntWithDefault("GRPC_PORT", 90)
	rPort := utils.GetEnvIntWithDefault("REST_PORT", 80)

	//	delay := flag.Uint("DELAY", uint(d), "delay for sending operations (ms)")
	grpcPort := flag.Uint("GRPC_PORT", uint(gPort), "port number of the grpc server")
	restPort := flag.Uint("REST_PORT", uint(rPort), "port number of the rest server")
	//	numThreads := flag.Uint("NUM_THREADS", uint(nt), "number of threads used to multicast messages")
	//	verb := flag.String("VERBOSE", verbose, "Turn verbose mode on or off.")
	registry_addr := flag.String("REGISTRY_ADDR", ":90", "service registry adress")
	r := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)
	var err error
	if *r {
		services = append(services, rgstr.Registration)
	}
	if *application {
		services = append(services, basic.RegisterService)
	}
	log.Println("start")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("Connecting grpc server..")
		err = StartServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		wg.Done()
	}()

	handler.GrpcPort = *grpcPort
	if *application {
		handler.RegClient, err = clientregistry.Connect(*registry_addr)
		router := gin.Default()
		routerGroup := router.Group(*registry_addr)
		routerGroup.GET("/groups", handler.GetGroups)
		routerGroup.POST("/groups", handler.CreateGroup)

		wg.Add(1)
		log.Println("Starting application")
		log.Println("http server started...")
		err := router.Run(fmt.Sprintf(":%d", *restPort))
		if err != nil {
			log.Println("Error in starting http server", err.Error())
		}
		if err != nil {
			log.Println("Error in connect client to registry ", err.Error())
		}
		wg.Done()
	}

	wgChan := make(chan bool)

	go func() {
		wg.Wait()
		wgChan <- true
	}()
	log.Println("App started")

	select {
	case <-wgChan:
	}

}

func StartServer(programAddress string, grpcServices ...func(grpc.ServiceRegistrar) error) error {

	lis, err := net.Listen("tcp", programAddress)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	for _, grpcService := range grpcServices {
		err = grpcService(s)
		if err != nil {
			return err
		}

	}

	if err = s.Serve(lis); err != nil {
		return err
	}
	log.Println("grpc server start")
	return nil
}

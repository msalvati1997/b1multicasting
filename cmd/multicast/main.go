package main

import (
	"flag"
	_ "flag"
	"fmt"
	"github.com/labstack/echo"
	"github.com/msalvati1997/b1multicasting/handler"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	_ "github.com/msalvati1997/b1multicasting/pkg/basic"
	basic "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	clientregistry "github.com/msalvati1997/b1multicasting/pkg/reg/client"
	registry "github.com/msalvati1997/b1multicasting/pkg/reg/server"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var grpcL net.Listener

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
	reg := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *reg {
		services = append(services, registry.Registration)
	}
	if *application {
		services = append(services, basic.RegisterService)
	}
	log.Println("start")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("Connecting grpc server..")
		err := StartServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		wg.Done()
	}()
	if *application {
		wg.Add(1)
		log.Println("Starting application")

		e := echo.New()
		api := e.Group("/api/v1", serverHeader)
		api.POST("/groups/:multicastId", handler.CreateGroup)
		api.GET("/groups", handler.GetGroups)

		go func() {
			log.Println("http server started...")
			err := e.Start(fmt.Sprintf(":%d", restPort))
			if err != nil {
				log.Println("Error in starting http server", err.Error())
			}
		}()
		var err error
		handler.Registryclient, err = clientregistry.Connect(*registry_addr)
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
func serverHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("x-version", "Test/v1.0")
		return next(c)
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

	return nil
}

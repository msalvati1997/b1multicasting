package main

import (
	"flag"
	"fmt"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	basic "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	rgstr "github.com/msalvati1997/b1multicasting/pkg/registryservice/server"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

// @title MULTICAST API
// @version 1.0
// @description This is a multicast API
// @contact.email salvatimartina97@gmail.com
//@host localhost:8080
// @BasePath /multicast/v1
func main() {

	delay := flag.Uint("DELAY", uint(utils.GetEnvIntWithDefault("DELAY", 0)), "delay for sending operations (ms)")
	grpcPort := flag.Uint("GRPC_PORT", uint(utils.GetEnvIntWithDefault("GRPC_PORT", 90)), "port number of the grpc server")
	restPort := flag.Uint("REST_PORT", uint(utils.GetEnvIntWithDefault("REST_PORT", 80)), "port number of the rest server")
	restPath := flag.String("restPath", utils.GetEnvStringWithDefault("REST_PATH", "/multicast/v1"), "path of the rest api")
	numThreads := flag.Uint("NUM_THREADS", uint(utils.GetEnvIntWithDefault("NUM_THREADS", 1)), "number of threads used to multicast messages")
	verb := flag.Bool("VERBOSE", utils.GetEnvBoolWithDefault("VERBOSE", true), "Turn verbose mode on or off.")
	registry_addr := flag.String("REGISTRY_ADDR", "registry:90", "service registry adress")
	r := flag.Bool("REGISTRY", utils.GetEnvBoolWithDefault("REGISTRY", false), "start multicast registry")
	application := flag.Bool("APP", utils.GetEnvBoolWithDefault("APP", false), "start multicast application")

	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)
	var err error

	if *application {
		log.Println("Adding basic communication service to gRPC server")
		services = append(services, basic.RegisterService)
	}
	if *r {
		log.Println("Adding multicast registry service to gRPC server")
		services = append(services, rgstr.Registration)
	}
	log.Println("start")
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		err = StartServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		wg.Done()
	}()
	if *application {

		wg.Add(1)
		go func() {
			err := multicastapp.Run(*grpcPort, *restPort, *registry_addr, *restPath, *numThreads, *delay, *verb)
			if err != nil {
				log.Println("Error in running applicatioon", err.Error())
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
	log.Printf("Grpc-Server started at %v", lis.Addr().String())

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

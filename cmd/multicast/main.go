package main

import (
	multicastApp "b1multicasting/internal/myapi"
	"b1multicasting/internal/utils"
	serverservice "b1multicasting/pkg/basic/server"
	serverregistry "b1multicasting/pkg/registry/server"
	"flag"
	_ "flag"
	"fmt"
	_ "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"sync"
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
	registry := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

	flag.Parse()

	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *registry {
		services = append(services, serverregistry.Registration)
	}
	if *application {
		services = append(services, serverservice.RegisterService)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		err := serverservice.RunServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		w.Done()
	}(&wg)

	if *application {
		wg.Add(1)
		go func() {
			err := multicastApp.Run(*grpcPort, *restPort, *registry_addr, *numThreads, *delay, *verb)
			if err != nil {
				return
			}
			wg.Done()
		}()
	}
}

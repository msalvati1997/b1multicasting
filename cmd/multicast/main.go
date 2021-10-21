package main

import (
	"flag"
	_ "flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/msalvati1997/b1multicasting/api"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	_ "github.com/msalvati1997/b1multicasting/pkg/basic"
	basic "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	clientregistry "github.com/msalvati1997/b1multicasting/pkg/registry/client"
	registry "github.com/msalvati1997/b1multicasting/pkg/registry/server"
	_ "github.com/sirupsen/logrus"
	"github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"log"
	"net/http"
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
	restPort := flag.Uint("REST_PORT", uint(rPort), "port number of the rest server")
	numThreads := flag.Uint("NUM_THREADS", uint(nt), "number of threads used to multicast messages")
	verb := flag.String("VERBOSE", verbose, "Turn verbose mode on or off.")
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
	var wg sync.WaitGroup
	wg.Add(2)
	go func(w *sync.WaitGroup) {
		err := utils.StartServer(fmt.Sprintf(":%d", *grpcPort), services...)
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
// @contact.email salvatimartina97@gmail.com
// @host localhost
// @BasePath /api
func Run(grpcP uint, restPort uint, registryAddr string, numThreads uint, dl uint, verbose string) error {
	var err error
	api.GrpcPort = grpcP
	api.Registryclient, err = clientregistry.Connect(registryAddr)
	if err != nil {
		log.Println("error", err)
		return err
	}
	newRouter := mux.NewRouter()
	newRouter.HandleFunc("/groups", api.GetGroups).Methods("GET")
	newRouter.HandleFunc("/groups", api.CreateGroup).Methods("PUT")
	//utils2.GoPool.Initialize(int(numThreads))
	// mount swagger API documentation
	newRouter.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", restPort), newRouter))
	return err
}

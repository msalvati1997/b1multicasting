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
	clientregistry "github.com/msalvati1997/b1multicasting/pkg/reg/client"
	registry "github.com/msalvati1997/b1multicasting/pkg/reg/server"

	_ "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
)

var grpcL net.Listener
var httpL net.Listener
var tcpm cmux.CMux

// @title Orders API
// @version 1.0
// @description This is a sample service for managing groups multicast
// @contact.email salvatimartina97@gmail.com
// @host localhost
// @BasePath /api
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
	// Create a listener at the desired port.
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *restPort))
	if err != nil {
		log.Fatal(err)
	}

	defer func(l net.Listener) {
		err := l.Close()

		if err != nil {

		}
	}(l)
	// Create a cmux object.
	tcpm = cmux.New(l)
	// Declare the match for different services required.
	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	grpcL = tcpm.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL = tcpm.Match(cmux.HTTP1Fast())
	var wg sync.WaitGroup
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		err := StartServer(services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		w.Done()
	}(&wg)

	wg.Add(1)
	if *application {
		api.GrpcPort = *grpcPort
		newRouter := mux.NewRouter()
		newRouter.HandleFunc("/groups", api.GetGroups).Methods("GET")
		newRouter.HandleFunc("/groups", api.CreateGroup).Methods("PUT")
		newRouter.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)
		go func() {
			err := http.Serve(httpL, newRouter)
			log.Println("http server started.")
			if err != nil {

			}
		}()
		// Start cmux serving.
		if err = tcpm.Serve(); !strings.Contains(err.Error(),
			"use of closed network connection") {
			log.Fatal(err)
		}

		log.Println("Server listening on ", restPort)

		api.Registryclient, err = clientregistry.Connect(*registry_addr)
		if err != nil {
			log.Println("error", err)
		}
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

func StartServer(grpcServices ...func(grpc.ServiceRegistrar) error) error {

	s := grpc.NewServer()
	for _, grpcService := range grpcServices {
		err := grpcService(s)
		if err != nil {
			return err
		}
	}

	if err := s.Serve(grpcL); err != nil {
		return err
	}

	return nil
}

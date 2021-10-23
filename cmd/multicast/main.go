package main

import (
	"flag"
	"fmt"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	basic "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	rgstr "github.com/msalvati1997/b1multicasting/pkg/registryservice/server"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var grpcL net.Listener
var M cmux.CMux

func main() {

	d := utils.GetEnvIntWithDefault("DELAY", 0)
	nt := utils.GetEnvIntWithDefault("NUM_THREADS", 1)
	verbose := utils.GetEnvBoolWithDefault("VERBOSE", true)
	rg := utils.GetEnvBoolWithDefault("REGISTRY", false)
	app := utils.GetEnvBoolWithDefault("APP", false)
	gPort := utils.GetEnvIntWithDefault("GRPC_PORT", 80)
	rPort := utils.GetEnvIntWithDefault("REST_PORT", 80)
	restP := utils.GetEnvStringWithDefault("REST_PATH", "/multicast/api/v1")
	delay := flag.Uint("DELAY", uint(d), "delay for sending operations (ms)")
	grpcPort := flag.Uint("GRPC_PORT", uint(gPort), "port number of the grpc server")
	restPort := flag.Uint("REST_PORT", uint(rPort), "port number of the rest server")
	restPath := flag.String("restPath", restP, "path of the rest api")
	numThreads := flag.Uint("NUM_THREADS", uint(nt), "number of threads used to multicast messages")
	verb := flag.Bool("VERBOSE", verbose, "Turn verbose mode on or off.")
	registry_addr := flag.String("REGISTRY_ADDR", ":90", "service registry adress")
	r := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

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
	l, err := net.Listen("tcp", ":80")
	M = cmux.New(l)
	go func() {
		grpcL = M.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
		err = StartServer(fmt.Sprintf(":%d", *restPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		wg.Done()
	}()
	if *application {

		wg.Add(1)
		go func() {
			err := multicastapp.Run(*grpcPort, *restPort, *registry_addr, *restPath, *numThreads, *delay, *verb, M)
			if err != nil {
				log.Println("Error in running applicatioon", err.Error())
				return
			}
			wg.Done()
		}()
	} else {
		err := M.Serve()
		if err != nil {
			return
		}
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

	//lis, err := net.Listen("tcp", programAddress)
	//if err != nil {
	//	return err
	//}

	s := grpc.NewServer()
	for _, grpcService := range grpcServices {
		err := grpcService(s)
		if err != nil {
			return err
		}

	}
	var err error
	go func() {
		err = s.Serve(grpcL)
		if err != nil {
			log.Println("Error ", err.Error())
		}
	}()
	return err
}

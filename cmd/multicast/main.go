package main

import (
	"flag"
	_ "flag"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/msalvati1997/b1multicasting/internal/app"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	serverservice "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	registry "github.com/msalvati1997/b1multicasting/pkg/registry/client"
	serverregistry "github.com/msalvati1997/b1multicasting/pkg/registry/server"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
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
	reg := flag.Bool("REGISTRY", rg, "start multicast registry")
	application := flag.Bool("APP", app, "start multicast application")

	flag.Parse()
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *reg {
		services = append(services, serverregistry.Registration)
	}
	if *application {
		services = append(services, serverservice.RegisterService)
	}
	log.Println("start")
	var wg sync.WaitGroup
	wg.Add(2)
	go func(w *sync.WaitGroup) {
		err := serverservice.RunServer(fmt.Sprintf(":%d", *grpcPort), services...)
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

// Run initializes and executes the Rest HTTP server
func Run(grpcP uint, restPort uint, registryAddr string, numThreads uint, dl uint, verbose string) error {
	var err error
	app.GrpcPort = grpcP
	app.RegistryClient, err = registry.Connect(registryAddr)
	if err != nil {
		return err
	}
	beego.BConfig.AppName = "MulticastApp"
	if verbose == "ON" {
		beego.BConfig.Log = beego.LogConfig{AccessLogs: true}
		beego.BConfig.Log.AccessLogs = true
	}

	beego.BConfig.RunMode = "dev"
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.BConfig.Listen.HTTPPort = int(restPort)
	utils2.GoPool.Initialize(int(numThreads))
	beego.Run()
	return nil
}

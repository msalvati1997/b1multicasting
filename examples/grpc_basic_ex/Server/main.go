package main

import (
	_ "flag"
	"log"
	server "b1multicasting/pkg/basic/server"
)

func main() {

	//Start the server
	err := server.RunServer(":8080", server.RegisterService)
	if err != nil {
		return
	}
	log.Println("Server connected..")
}

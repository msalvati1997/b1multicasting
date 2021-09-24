package main

import (
	server "b1multicasting/pkg/basic/server"
	_ "flag"
	"log"
)

func main() {

	//Start the server
	err := server.RunServer(":8080", server.RegisterService)
	if err != nil {
		return
	}
	log.Println("Server connected..")
}

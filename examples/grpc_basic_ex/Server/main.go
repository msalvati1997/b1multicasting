package main

import (
	_ "flag"
	server "github.com/msalvati1997/b1multicasting/pkg/basic/server"
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

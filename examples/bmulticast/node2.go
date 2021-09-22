package main

import (
	"b1multicasting/pkg/basic"
	server "b1multicasting/pkg/basic/server"
	"b1multicasting/pkg/multicasting"
	"bufio"
	"log"
	"os"
)

func main() {
	//effettuo il run del server in una go-routines
	go func() {
		err := server.RunServer(":8081", server.RegisterService)
		if err != nil {
			return
		}
	}()
	//effettuo la connessione degli altri nodi come clients
	nodesPort := []string{":8082", ":8080"}
	delay := 1
	multicasterId := "MulticasterId"
	connections, err := multicasting.Connections(nodesPort, delay)
	if err != nil {
		log.Println("Error in connecting Clients ", err.Error())
	}
	for {
		log.Println("Input : ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			log.Println("Input : ")
			text := scanner.Bytes()
			err := connections.XMulticast(multicasterId, basic.NewMessage(text))
			if err != nil {
				log.Println(err)
			}
		}
	}

}

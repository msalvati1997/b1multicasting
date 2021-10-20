package main

import (
	"bufio"
	"flag"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	server "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	"os"
	"strings"
)

func main() {
	port := flag.String("port", ":8080", "port number of the server")
	membersPort := flag.String("membersPort", ":8081,:8082", "ports of the member of the multicast group")
	multicasterId := flag.String("multicastId", "MulticasterId", "id of the multicaster id")
	delay := flag.Int("delay", 0, "delay of sending operation")
	flag.Parse()
	go func() {
		err := server.RunServer(*port, server.RegisterService)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
	}()
	//effettuo la connessione degli altri nodi come clients
	member := strings.Split(*membersPort, ",")
	member = append(member, *port)
	connections, err := multicasting.Connections(member, *delay)
	if err != nil {
		log.Println("Error in connecting Clients ", err.Error())
	}
	log.Println("Input : ")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Bytes()
			msg := basic.NewMessage(make(map[string]string), text)
			err := connections.BMulticast(*multicasterId, msg)
			if err != nil {
				log.Println(err)
			}
			log.Println("Input : ")
		}
	}
}

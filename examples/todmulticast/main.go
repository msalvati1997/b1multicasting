package main

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	server "b1multicasting/pkg/basic/server"
	"b1multicasting/pkg/multicasting"
	utils2 "b1multicasting/pkg/utils"
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	//effettuo il run del server in una go-routines
	port := flag.String("port", ":8080", "port number of the server")
	membersPort := flag.String("membersPort", ":8081,:8082", "ports of the member of the multicast group")
	multicasterId := flag.String("multicastId", "MulticasterId", "id of the multicaster id")
	delay := flag.Int("delay", 0, "delay of sending operation")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
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
	Connections, _ := multicasting.Connections(member, *delay)
	myport, _ := strconv.Atoi(strings.Split(*port, ":")[1])
	VectorId := myport
	numberOfThreads := 10
	utils2.GoPool.Initialize(numberOfThreads, Connections)

	log.Println("Input : ")
	scanner := bufio.NewScanner(os.Stdin)
	go utils2.TODDeliver()

	for scanner.Scan() {
		text := scanner.Bytes()
		msg := basic.NewMessage(make(map[string]string), text)
		id := utils.GenerateUID()
		msg.MessageHeader["ProcessId"] = strings.Split(*port, ":")[1]
		msg.MessageHeader["i"] = id
		msg.MessageHeader["s"] = strconv.FormatUint(utils.Clock.Tock(), 10)
		msg.MessageHeader["type"] = "TOD"
		msg.MessageHeader["ProcessId"] = strconv.Itoa(VectorId)
		msg.MessageHeader["GroupId"] = *multicasterId
		utils2.GoPool.MessageCh <- msg
		log.Println("Input : ")
	}
}

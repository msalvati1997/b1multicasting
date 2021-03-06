package main

import (
	"bufio"
	"flag"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	server "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"log"
	"os"
	"strings"
)

func main() {
	port := flag.String("port", ":8080", "port number of the server")
	membersPort := flag.String("membersPort", ":8081,:8082", "ports of the member of the multicast group")
	multicasterId := flag.String("multicastId", "MulticastId", "id of the multicast group")
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
	Connections, err := multicasting.Connections(member, *delay)
	if err != nil {
		log.Println("Error in connecting Clients ", err.Error())
		return
	}
	log.Println("Input : ")
	//selection of the sequencer
	//the sequencer is one of the nodes partecipating in multicasting
	sequencerPort := multicasting.SelectingSequencer(member, false)
	seqCon, err := Connections.GetGrpcClient(sequencerPort)
	if err != nil {
		log.Println("Error in find connection with sequencer..", err.Error())
		return
	}
	seq := false
	if *port == sequencerPort {
		log.Println("I'm the sequencer of MulticastGroup", *multicasterId)
		seq = true
	} else {
		log.Println("The sequencer nodes is at port", sequencerPort)
	}

	multicasting.Seq.MulticastId = *multicasterId
	multicasting.Seq.SeqConn = *seqCon
	multicasting.Seq.Conns = *Connections
	multicasting.Seq.B = seq
	multicasting.Seq.SeqPort = sequencerPort
	numberOfThreads := 10
	utils2.GoPool.Initialize(numberOfThreads)
	go utils2.TOCDeliver()
	for {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Bytes()
			msg := basic.NewMessage(make(map[string]string), text)
			msg.MessageHeader["type"] = "TOC"
			msg.MessageHeader["GroupId"] = *multicasterId
			//Sender attaches the unique id to the message and sends <m,i> to the sequencer as well as to the group
			utils2.GoPool.MessageCh <- msg
			log.Println("Input : ")
		}
	}
}

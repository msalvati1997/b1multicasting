package main

import (
	"bufio"
	"flag"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	server "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
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
	member := strings.Split(*membersPort, ",")
	member = append(member, *port)
	myport, _ := strconv.Atoi(strings.Split(*port, ":")[1])
	_, err := multicasting.Connections(member, *delay)
	if err != nil {
		return
	}
	IndexProcesses := make([]int, 0, len(member))
	for i := 0; i < len(member); i++ {
		m := strings.Split(member[i], ":")[1]
		mi, _ := strconv.Atoi(m)
		IndexProcesses = append(IndexProcesses, mi)
	}
	sort.Ints(IndexProcesses)
	mapProcessId := make(map[int]int, len(member))
	for i := 0; i < len(member); i++ {
		if IndexProcesses[i] == myport {
			utils.Myid = i
		}
		mapProcessId[IndexProcesses[i]] = i
	}
	log.Println(utils.Myid)
	log.Println("Input : ")
	utils.Vectorclock = utils.NewVectorClock(len(member))
	numberOfThreads := 10
	utils2.GoPool.Initialize(numberOfThreads)
	//VectorId := myport
	scanner := bufio.NewScanner(os.Stdin)
	go utils2.CODeliver()

	for scanner.Scan() {
		text := scanner.Bytes()
		msg := basic.NewMessage(make(map[string]string), text)
		id := utils.GenerateUID()
		msg.MessageHeader["ProcessId"] = strings.Split(*port, ":")[1]
		msg.MessageHeader["i"] = id
		msg.MessageHeader["type"] = "CO"
		msg.MessageHeader["GroupId"] = *multicasterId
		utils2.GoPool.MessageCh <- msg
		log.Println("Input : ")
	}
}

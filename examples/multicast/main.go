package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	server "github.com/msalvati1997/b1multicasting/pkg/basic/server"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice"
	clientregistry "github.com/msalvati1997/b1multicasting/pkg/registryservice/client"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	registry "github.com/msalvati1997/b1multicasting/pkg/registryservice/server"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	grpcPort := flag.Uint("grpcPort", 90, "port number of the grpc server")
	app := flag.Bool("appMulticast", false, "start distributed multicast node")
	registryMulticast := flag.Bool("registryMulticast", false, "start multicast registry")
	registryAddr := flag.String("registryAddr", ":90", "multicast registry address")
	numThreads := flag.Uint("numThreads", 1, "number of threads used to multicast messages")
	multicastId := flag.String("multicastId", "multicastTopic", "identifier of the multicast group")
	delay := flag.Uint("delay", 0, "delay for sending operations (ms)")
	startMulticastGroup := flag.Bool("startMulticastGroup", false, "start multicast group after registering")
	multicastTypeString := flag.String("multicastType", "BMULTICAST", fmt.Sprintf("multicast typology: %s", registryservice.MulticastTypes))
	flag.Parse()

	multicastType, ok := registryservice.MulticastType[*multicastTypeString]

	if !ok {
		log.Println("Multicast type %s not supported", *multicastTypeString)
	}
	services := make([]func(registrar grpc.ServiceRegistrar) error, 0)

	if *registryMulticast {
		services = append(services, registry.Registration)
	}
	if *app {
		services = append(services, server.RegisterService)
	}
	log.Println("start")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("Connecting grpc server ", *grpcPort)
		err := StartServer(fmt.Sprintf(":%d", *grpcPort), services...)
		if err != nil {
			log.Println("Error in connecting server", err.Error())
			return
		}
		wg.Done()
	}()
	if *registryMulticast {
		for {

		}
	}
	if *app {
		log.Println("Multicast type number", multicastType.Number())
		utils2.GoPool.Initialize(int(*numThreads))
		registryClient, err := clientregistry.Connect(*registryAddr)

		if err != nil {
			log.Println("error in connecting to registry ", err.Error())
			return
		}
		log.Println("Connected to registry")

		registrationAns, err := registryClient.Register(context.Background(), &protoregistry.Rinfo{
			MulticastId:   *multicastId,
			MulticastType: multicastType,
			ClientPort:    uint32(*grpcPort),
		})

		if err != nil {
			log.Println("Error in register to a group", err.Error())
			return
		}
		log.Println("Register to group")

		groupInfo := registrationAns.GroupInfo
		clientId := registrationAns.ClientId

		if *startMulticastGroup {
			groupInfo, err = registryClient.StartGroup(context.Background(), &protoregistry.RequestData{
				MulticastId: *multicastId,
				MId:         clientId})
			if err != nil {
				log.Println("Error in start group ", err.Error())
				return
			}
		}
		for groupInfo.Status == protoregistry.Status_OPENING {
			time.Sleep(time.Second * 5)
			groupInfo, err = registryClient.GetStatus(context.Background(), &protoregistry.MulticastId{MulticastId: *multicastId})
			if err != nil {
				log.Println("Error in get status ", err.Error())
			}

		}
		if groupInfo.Status == protoregistry.Status_CLOSED {
			log.Println("Error in group closed ", err.Error())
		}

		groupInfo, err = registryClient.Ready(context.Background(), &protoregistry.RequestData{MulticastId: *multicastId, MId: clientId})
		log.Println("Ready")
		if err != nil {
			log.Println("Error in Ready group ", err.Error())

		}
		for groupInfo.Status == protoregistry.Status_STARTING {
			time.Sleep(time.Second * 5)
			groupInfo, err = registryClient.GetStatus(context.Background(), &protoregistry.MulticastId{MulticastId: *multicastId})
			if err != nil {
				log.Println("Error in get status ", err.Error())

			}
		}
		if groupInfo.Status == protoregistry.Status_CLOSED {
			log.Println("Group closed")
		}
		var members []string
		//effettuo la connessione degli altri nodi come clients
		for memberId, member := range groupInfo.Members {
			if memberId != clientId {
				log.Println("Connecting to: %s", member.Address)
				members = append(members, member.Address)
			}
		}
		members = append(members, clientId)
		Connections, err := multicasting.Connections(members, int(*delay))
		if err != nil {
			log.Println("Error in connecting Clients ", err.Error())
			return
		}
		if multicastType.Number() == 0 {
			log.Println("STARTING BMULTICAST COMMUNICATION")
			log.Println("Input : ")
			for {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					text := scanner.Bytes()
					msg := basic.NewMessage(make(map[string]string), text)
					err := Connections.BMulticast(*multicastId, msg)
					if err != nil {
						log.Println(err)
					}
					log.Println("Input : ")
				}
			}
		}
		if multicastType.Number() == 1 {
			log.Println("STARTING TOC COMMUNICATION")
			sequencerPort := multicasting.SelectingSequencer(members, false)
			seqCon, err := Connections.GetGrpcClient(sequencerPort)
			if err != nil {
				log.Println("Error in find connection with sequencer..", err.Error())
				return
			}
			seq := false
			if clientId == sequencerPort {
				log.Println("I'm the sequencer of MulticastGroup", *multicastId)
				seq = true
			} else {
				log.Println("The sequencer nodes is at port", sequencerPort)
			}
			multicasting.Seq.MulticastId = *multicastId
			multicasting.Seq.SeqConn = *seqCon
			multicasting.Seq.Conns = *Connections
			multicasting.Seq.B = seq
			multicasting.Seq.SeqPort = sequencerPort
			go utils2.TOCDeliver()
			for {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					text := scanner.Bytes()
					msg := basic.NewMessage(make(map[string]string), text)
					msg.MessageHeader["i"] = utils.GenerateUID()
					msg.MessageHeader["type"] = "TOC"
					msg.MessageHeader["GroupId"] = *multicastId
					//Sender attaches the unique id to the message and sends <m,i> to the sequencer as well as to the group
					utils2.GoPool.MessageCh <- msg
					log.Println("Input : ")
				}
			}
		}
		if multicastType.Number() == 2 {
			log.Println("STARTING TOD COMMUNICATION")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				text := scanner.Bytes()
				msg := basic.NewMessage(make(map[string]string), text)
				id := utils.GenerateUID()
				msg.MessageHeader["i"] = id
				msg.MessageHeader["s"] = strconv.FormatUint(utils.Clock.Tock(), 10)
				msg.MessageHeader["type"] = "TOD"
				msg.MessageHeader["ProcessId"] = clientId
				msg.MessageHeader["GroupId"] = *multicastId
				utils2.GoPool.MessageCh <- msg
				log.Println("Input : ")
			}
			go utils2.TODDeliver()
		}
		if multicastType.Number() == 3 {
			utils.Vectorclock = utils.NewVectorClock(len(groupInfo.Members))
			log.Println("STARTING CO COMMUNICATION")
			scanner := bufio.NewScanner(os.Stdin)
			go utils2.CODeliver()
			for scanner.Scan() {
				text := scanner.Bytes()
				msg := basic.NewMessage(make(map[string]string), text)
				id := utils.GenerateUID()
				msg.MessageHeader["ProcessId"] = clientId
				msg.MessageHeader["i"] = id
				msg.MessageHeader["type"] = "CO"
				msg.MessageHeader["GroupId"] = *multicastId
				utils2.GoPool.MessageCh <- msg
				log.Println("Input : ")
			}
		}
	}

}

func StartServer(programAddress string, grpcServices ...func(grpc.ServiceRegistrar) error) error {

	lis, err := net.Listen("tcp", programAddress)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	for _, grpcService := range grpcServices {
		err = grpcService(s)
		if err != nil {
			return err
		}

	}

	if err = s.Serve(lis); err != nil {
		return err
	}

	return nil
}

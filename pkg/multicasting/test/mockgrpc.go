package test

import (
	"context"
	utils2 "github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	pb "github.com/msalvati1997/b1multicasting/pkg/basic/proto"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"math"
	"net"
	"sort"
	"strconv"
)

type mockServer struct {
	pb.UnimplementedEndToEndServiceServer
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterEndToEndServiceServer(server, &mockServer{})
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func ClientConnMock(target string, delay int) (*client.GrpcClient, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	cl := new(client.Client)
	cl.Client = pb.NewEndToEndServiceClient(conn)
	cl.Connection = conn
	cl.Address = target
	cl.Delay = delay
	return &client.GrpcClient{cl}, nil
}

// implementationof the service methods called by Grpc
func (s *mockServer) SendMessage(ctx context.Context, in *pb.RequestMessage) (*pb.ResponseMessage, error) {
	//source, _ := peer.FromContext(ctx)
	id := in.GetId()
	log.Println("Request to ", in.GetId(), " \t Received message : payload: ", string(in.Payload), " header ", in.MessageHeader)
	if in.MessageHeader["type"] == "TOD" {
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		msg.MessageHeader["ReceiverId"] = id
		i, _ := strconv.Atoi(id)
		BDeliverTOD(msg, i)
	}
	if in.MessageHeader["type"] == "ACK" {
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		index := in.MessageHeader["i"]
		seq := in.MessageHeader["s"]
		nseq, _ := strconv.Atoi(seq)
		node := in.MessageHeader["nd"]
		log.Println("Received ack from ", node, " to ", id)
		n, _ := strconv.Atoi(id)
		nodo := GetNode(n)
		nodo.MuAck.Lock()
		nodo.AckQueue.Q = append(nodo.AckQueue.Q, multicasting.SeqMessage{
			Msg:  msg,
			I:    index,
			Nseq: nseq,
		})
		if len(nodo.AckQueue.Q) > 1 {
			sort.Slice(nodo.AckQueue.Q, func(i, j int) bool {
				return nodo.AckQueue.Q[i].Nseq <= nodo.AckQueue.Q[j].Nseq
			})
		}
		utils.PrintACKQueue_(&nodo.AckQueue, n)
		nodo.MuAck.Unlock()
	}
	if in.MessageHeader["type"] == "TOC" && in.MessageHeader["seq"] == "true" {
		in.MessageHeader["delnode"] = id
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		n, _ := strconv.Atoi(id)
		BDeliverSeq(msg, n)
	}
	if in.MessageHeader["order"] == "true" && in.MessageHeader["type"] == "TOC2" {
		n, _ := strconv.Atoi(id)
		TOCDeliver(basic.NewMessage(in.MessageHeader, in.Payload), n)
	}
	if in.MessageHeader["type"] == "CO" {
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		n, _ := strconv.Atoi(id)
		myn, _ := strconv.Atoi(in.MessageHeader["ProcessId"])
		if n != myn { //il processo ha già deliverato il messaggio che ha inviato in multicast
			BDeliverCO(msg, n)
		}
	}
	return &pb.ResponseMessage{}, nil
}

func BDeliverCO(msg basic.Message, id int) {
	itsvector := utils2.NewVectorClock(multicasting.GetNumbersOfClients())
	for i := 0; i < multicasting.GetNumbersOfClients(); i++ {
		p := msg.MessageHeader[strconv.Itoa(i)]
		d, _ := strconv.Atoi(p)
		itsvector.LeapI(i, uint64(d))
	}
	//place <Vig,m> in holdbackqueue
	AppendMessageToCOQueue(msg, itsvector, id)
}

func AppendMessageToCOQueue(msg basic.Message, v utils2.VectorClock, id int) {
	nodo := GetNode_co(id)
	nodo.MuMsg.Lock()
	nodo.HoldBackQueue.Q = append(nodo.HoldBackQueue.Q, multicasting.COMessage{
		Msg:    msg,
		Vector: v,
	})
	nodo.MuMsg.Unlock()
}

func TOCDeliver(msg basic.Message, id int) {
	log.Println("TOCDELIVER of message ", string(msg.Payload), " of ", id, "with seq ", msg.MessageHeader["s"])
	seq := msg.MessageHeader["s"]
	s, _ := strconv.Atoi(seq)
	AppendMessageTct(msg, msg.MessageHeader["i"], uint64(s), id)
}

func BDeliverSeq(msg basic.Message, id int) {
	M.Lock()
	log.Println("Delivering by sequencer ", id, "of message ", string(msg.Payload)+" to other member with Sg", Sg)
	msg.MessageHeader["s"] = strconv.Itoa(Sg)
	msg.MessageHeader["order"] = "true"
	msg.MessageHeader["type"] = "TOC2"
	AppendMessageTct(msg, msg.MessageHeader["i"], uint64(Sg), id)
	m := basic.NewMessage(msg.MessageHeader, msg.Payload)
	err := Node2_toc.SeqStruct.Conns.BMulticast(Node2_toc.SeqStruct.MulticastId, m)
	Sg = Sg + 1 //update Sg = Sg+1
	log.Println(Sg)
	if err != nil {
		return
	}
	M.Unlock()
}

func BDeliverTOD(message basic.Message, n int) {
	//msgi viene posto da ogni processo ricevente pj¬ in una coda locale queuej, ordinata in base al valore del timestamp.
	log.Println("Delivering of ", n, "of the message ", message.MessageHeader["i"], "receveid from ", message.MessageHeader["ProcessId"])
	id, _ := message.MessageHeader["i"]
	cl, _ := strconv.Atoi(message.MessageHeader["s"])
	nodo := GetNode(n)
	nodo.MuClock.Lock()
	max := math.Max(float64(nodo.Clock.Tock()), float64(cl))
	nodo.Clock.Leap(uint64(max) + 1)
	nodo.MuMsg.Lock()
	AppendMessageTOD(message, id, uint64(cl), n) //tutti i processi aggiungono all'holdbackqueue il messaggio con il massimo numero di sequenza
	//Il processo ricevente invia in multicast un messaggio di Ack con il max # seq
	m := make(map[string]string)
	m["type"] = "ACK"
	m["i"] = message.MessageHeader["i"]
	m["nd"] = strconv.Itoa(n)
	m["s"] = strconv.Itoa(cl)
	m["GroupId"] = "TODMULTICAST_TEST"
	//ACKMulticast
	go func() {
		err := multicasting.Cnn.ACKMulticast(m["GroupId"], basic.NewMessage(m, []byte("ACK"+m["nd"])))
		if err != nil {
			return
		}
	}()
	nodo.MuClock.Unlock()
}

func MockConnections(ports []string, delay int) (*multicasting.Conns, error) {
	clients := make([]*client.GrpcClient, len(ports))
	for i := 0; i < len(ports); i++ {
		log.Println("Connecting with", ports[i], ",with max delay", delay)
		conn, err := ClientConnMock(ports[i], delay)
		if err != nil {
			log.Println("Error in connecting Clients ", err.Error())
			return nil, err
		}
		clients[i] = conn
	}
	multicasting.Cnn.Conns = clients
	return &multicasting.Conns{Conns: clients}, nil
}

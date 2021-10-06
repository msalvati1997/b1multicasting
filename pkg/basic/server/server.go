package basic

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/basic/proto"
	"b1multicasting/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"strconv"
)

type Server struct {
	proto.UnimplementedEndToEndServiceServer
}

//method that register the Grpc Server service
func RegisterService(s grpc.ServiceRegistrar) (err error) {
	proto.RegisterEndToEndServiceServer(s, &Server{})
	return
}

//implementation of the service methods called by Grpc
func (s *Server) SendMessage(ctx context.Context, in *proto.RequestMessage) (*proto.ResponseMessage, error) {
	source, _ := peer.FromContext(ctx)
	id := in.GetId()
	log.Println("Request from :{user_ip :", source.Addr, ",id : ", id, "} ")
	log.Println("Received message : payload: ", string(in.Payload), " header ", in.MessageHeader)
	//start deliverying..
	//if in.MessageHeader["type"]=="TOC" && in.MessageHeader["seq"]=="true" {
	//	in.MessageHeader["delnode"] = id
	//	node := utils.DelivererNode{NodeId: id}
	//	msg := basic.NewMessage(in.MessageHeader, in.Payload)
	//	node.BDeliverSeq(msg)
	//}
	if in.MessageHeader["type"] == "TOC" && in.MessageHeader["member"] == "false" && in.MessageHeader["seq"] == "true" {
		in.MessageHeader["delnode"] = id
		node := utils.DelivererNode{NodeId: id}
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		node.BDeliverSeq(msg)
	}
	if in.MessageHeader["type"] == "TOC" && in.MessageHeader["member"] == "true" && in.MessageHeader["seq"] == "false" {
		in.MessageHeader["delnode"] = id
		node := utils.DelivererNode{NodeId: id}
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		node.BDeliverMember(msg)
	}
	if in.MessageHeader["type"] == "Basic" {
		node := utils.DelivererNode{NodeId: id}
		node.BDeliver(basic.NewMessage(in.MessageHeader, in.Payload))
	}
	if in.MessageHeader["order"] == "true" && in.MessageHeader["type"] == "TOC2" {
		node := utils.DelivererNode{NodeId: id}
		node.TOCDeliver(basic.NewMessage(in.MessageHeader, in.Payload))
	}
	if in.MessageHeader["type"] == "TOD" {
		node := utils.DelivererNode{NodeId: id}
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		msg.MessageHeader["delnode"] = id
		node.BDeliverTOD(msg)
	}
	if in.MessageHeader["type"] == "ACK" {
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		index := in.MessageHeader["i"]
		seq := in.MessageHeader["s"]
		nseq, _ := strconv.Atoi(seq)
		utils.AppendACK(msg, index, nseq)
	}
	if in.MessageHeader["type"] == "CO" {
		msg := basic.NewMessage(in.MessageHeader, in.Payload)
		node := utils.DelivererNode{NodeId: id}
		msg.MessageHeader["delnode"] = id
		node.BDeliverCO(msg)
	}
	return &proto.ResponseMessage{}, nil
}

func RunServer(programAddress string, grpcServices ...func(grpc.ServiceRegistrar) error) error {
	//listening over the port
	n, err := getNetListener(programAddress)
	if err != nil {
		log.Println("Error in listening at port", programAddress)
		return err
	}
	//start new grpc server
	s := grpc.NewServer()

	//register the grpc service over the server
	for _, grpcService := range grpcServices {
		err = grpcService(s)
		if err != nil {
			return err
		}
	}
	//server
	log.Println("server connected")
	err = s.Serve(n)
	if err != nil {
		log.Println("failed to serve: %s", err)
		return err
	}
	return nil
}

func getNetListener(port string) (net.Listener, error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return lis, err
}

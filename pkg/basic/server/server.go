package basic

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/basic/proto"
	"b1multicasting/pkg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"net"
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
	//start deliverying..
	log.Println("Processing message..")
	node := utils.DelivererNode{NodeId: id}
	node.BDeliver(basic.NewMessage(in.Payload))
	log.Println("Message processed")
	return &proto.ResponseMessage{}, nil
}

func RunServer(programAddress string, grpcServices ...func(grpc.ServiceRegistrar) error) error {
	//listening over the port
	net, err := getNetListener(programAddress)
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
	err = s.Serve(net)
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
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	return lis, err
}

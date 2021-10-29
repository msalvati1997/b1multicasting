package basic

import (
	utils2 "github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/basic/proto"
	"github.com/msalvati1997/b1multicasting/pkg/multicastapp"
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"github.com/msalvati1997/b1multicasting/pkg/utils"
	"golang.org/x/net/context"
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

	if in.MessageHeader["Tranport"] == "HTTP" {
		mid := in.MessageHeader["GroupId"]
		group := multicastapp.MulticastGroups[mid]
		group.Group.ReceivedMessages = group.Group.ReceivedMessages + 1
		group.MessageMu.Lock()
		defer group.MessageMu.Unlock()
		msgh := multicastapp.Message{
			MessageHeader: in.MessageHeader,
			Payload:       in.Payload,
		}
		group.Messages = append(group.Messages, msgh)
	}
	source, _ := peer.FromContext(ctx)
	id := in.GetId()
	log.Println("Request from :{user_ip :", source.Addr, ",auth : ", source.AuthInfo, "} ")
	log.Println("Received message : payload: ", string(in.Payload), " header ", in.MessageHeader)
	//start deliverying..
	if in.MessageHeader["type"] == "TOC" && in.MessageHeader["seq"] == "true" {
		in.MessageHeader["delnode"] = id
		node := utils.DelivererNode{NodeId: id}
		node.BDeliverSeq(basic.NewMessage(in.MessageHeader, in.Payload))
	}
	//if in.MessageHeader["type"] == "TOC" && in.MessageHeader["member"] == "false" && in.MessageHeader["seq"] == "true" {
	//	in.MessageHeader["delnode"] = id
	//	node := utils.DelivererNode{NodeId: id}
	//	msg := basic.NewMessage(in.MessageHeader, in.Payload)
	//	node.BDeliverSeq(msg)
	//}
	//if in.MessageHeader["type"] == "TOC" && in.MessageHeader["member"] == "true" && in.MessageHeader["seq"] == "false" {
	//	in.MessageHeader["delnode"] = id
	//	node := utils.DelivererNode{NodeId: id}
	//	msg := basic.NewMessage(in.MessageHeader, in.Payload)
	//	node.BDeliverMember(msg)
	//}
	if in.MessageHeader["type"] == "B" {
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
		sender, _ := strconv.Atoi(in.MessageHeader["ProcessId"])
		if sender != utils2.Myid { //il processo ha già deliverato il messaggio che ha inviato in multicast
			node.BDeliverCO(msg)
		}
	}
	if in.MessageHeader["type"] == "CloseGroup" {
		groupInfo, err := multicastapp.RegClient.CloseGroup(context.Background(), &protoregistry.RequestData{
			MulticastId: in.MessageHeader["MulticastId"],
			MId:         id,
		})
		if err != nil {
			log.Println("Error in closing group")
		}
		multicastapp.MulticastGroups[groupInfo.MulticastId].Group.Status = protoregistry.Status_CLOSED.String()
		for key, _ := range multicastapp.MulticastGroups[groupInfo.MulticastId].Group.Members {
			member1 := multicastapp.MulticastGroups[groupInfo.MulticastId].Group.Members[key]
			member1.Ready = false
			multicastapp.MulticastGroups[groupInfo.MulticastId].Group.Members[key] = member1
		}
		for i := 0; i < len(utils.Del.DelivererNodes); i++ {
			if utils.Del.DelivererNodes[i].M.MessageHeader["GroupId"] == in.MessageHeader["MulticastId"] {
				utils.Del.DelivererNodes = append(utils.Del.DelivererNodes[:i], utils.Del.DelivererNodes[i+1:]...)
			}
		}
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

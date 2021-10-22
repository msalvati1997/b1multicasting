package client

import (
	"github.com/msalvati1997/b1multicasting/pkg/reg/proto"
	"google.golang.org/grpc"
	"log"
)

func Connect(address string) (proto.RegistryClient, error) {

	log.Println("Connecting to registry server ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("Error in connecting to registry server to adress ", address, " ", err.Error())
		return nil, err
	}

	return proto.NewRegistryClient(conn), nil
}

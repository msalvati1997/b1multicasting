package client

import (
	"github.com/msalvati1997/b1multicasting/pkg/reg/proto"
	"google.golang.org/grpc"
	"log"
)

func Connect(address string) (proto.RegistryClient, error) {

	log.Println("Connecting to registry server")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return proto.NewRegistryClient(conn), nil
}

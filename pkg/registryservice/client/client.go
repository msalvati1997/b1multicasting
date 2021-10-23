package client

import (
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
	"google.golang.org/grpc"
	"log"
)

func Connect(address string) (protoregistry.RegistryClient, error) {

	log.Println("Connecting to registry server ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println("Error in connecting to registry server to address ", address, " ", err.Error())
		return nil, err
	} else {
		log.Println("Correctly connected to registry server ")
	}

	return protoregistry.NewRegistryClient(conn), nil
}

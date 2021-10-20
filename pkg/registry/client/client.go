package client

import (
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	"google.golang.org/grpc"
)

func Connect(address string) (proto.RegistryClient, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return proto.NewRegistryClient(conn), nil
}

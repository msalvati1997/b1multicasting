package utils

import (
	"google.golang.org/grpc"
	"net"
)

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

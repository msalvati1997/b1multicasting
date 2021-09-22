package multicasting

import (
	client "b1multicasting/pkg/basic/client"
	"log"
)

type Conns struct {
	conns []*client.GrpcClient
}

func Connections(ports []string, delay int) (*Conns, error) {

	clients := make([]*client.GrpcClient, len(ports))
	for i := 0; i < len(ports); i++ {
		log.Println("Connecting with", ports[i])
		conn, err := client.Connect(ports[i], uint(delay))
		if err != nil {
			log.Println("Error in connecting Clients ", err.Error())
			return nil, err
		}
		clients[i] = conn
	}
	return &Conns{clients}, nil
}

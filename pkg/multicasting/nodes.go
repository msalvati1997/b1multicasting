package multicasting

import (
	client "b1multicasting/pkg/basic/client"
	"log"
)

type Conns struct {
	conns []*client.GrpcClient
}

var (
	Cnn Conns
)

func GetNumbersOfClients() int {
	return len(Cnn.conns)
}

func init() {
	Cnn.conns = make([]*client.GrpcClient, 0, 100)
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
	Cnn.conns = clients
	return &Conns{clients}, nil
}

func GetConns() []*client.GrpcClient {
	return Cnn.conns
}

func (c *Conns) GetGrpcClient(target string) (*client.GrpcClient, error) {
	for i := 0; i < len(c.conns); i++ {
		if c.conns[i].GetTarget() == target {
			return c.conns[i], nil
		}
	}
	return &client.GrpcClient{}, nil
}

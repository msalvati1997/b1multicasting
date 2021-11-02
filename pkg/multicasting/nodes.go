package multicasting

import (
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	"log"
)

type Conns struct {
	Conns []*client.GrpcClient
}

var (
	Cnn Conns
)

func GetNumbersOfClients() int {
	return len(Cnn.Conns)
}

func init() {
	Cnn.Conns = make([]*client.GrpcClient, 0, 100)
}

//method that permits multiple grpc connections from addresses
func Connections(addresses []string, delay int) (*Conns, error) {

	clients := make([]*client.GrpcClient, len(addresses))
	for i := 0; i < len(addresses); i++ {
		log.Println("Connecting with", addresses[i])
		conn, err := client.Connect(addresses[i], delay)
		if err != nil {
			log.Println("Error in connecting Clients ", err.Error())
			return nil, err
		}
		clients[i] = conn
	}
	Cnn.Conns = clients
	return &Conns{clients}, nil
}

func GetConns() []*client.GrpcClient {
	return Cnn.Conns
}

func (c *Conns) GetGrpcClient(target string) (*client.GrpcClient, error) {
	for i := 0; i < len(c.Conns); i++ {
		if c.Conns[i].GetTarget() == target {
			return c.Conns[i], nil
		}
	}
	return &client.GrpcClient{}, nil
}

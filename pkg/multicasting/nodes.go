package multicasting

import (
	"b1multicasting/pkg/basic"
	client "b1multicasting/pkg/basic/client"
	"log"
)

type Conns struct {
	conns []*client.GrpcClient
}

func (c *Conns) Multicast(id string, message basic.Message) error {

	for i := 0; i < len(c.conns); i++ {
		err := c.conns[i].Send(id, message)
		if err != nil {
			return err
		}
	}
	return nil
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

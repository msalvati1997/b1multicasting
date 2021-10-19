package basic

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/basic/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Client struct {
	Client     proto.EndToEndServiceClient
	Connection *grpc.ClientConn
	Delay      int //millisecond
	Address    string
}

type GrpcClient struct {
	Client *Client
}

//method to connect to grpcServer
func Connect(address string, delay int) (*GrpcClient, error) {
	log.Println("Connecting with server ", address)
	cc, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("connection state ====> ", cc.GetState(), "connected client ", cc.Target())

	cl := new(Client)
	cl.Client = proto.NewEndToEndServiceClient(cc)
	cl.Connection = cc
	cl.Address = address
	cl.Delay = delay

	return &GrpcClient{cl}, nil
}

//method to send message to GrpcServer
func (c *GrpcClient) Send(id string, message basic.Message, ch *chan bool) error {
	var w sync.WaitGroup

	if c.Client == nil {
		panic("Closed Connection")
	}
	if c.Client.Delay > 0 {
		max := c.Client.Delay * 1000
		min := 10
		v := rand.Intn(max-min) + min
		w.Add(1)
		WaitDelay(v, &w)
	}
	defer w.Wait()
	_, err := c.Client.Client.SendMessage(context.Background(),
		&proto.RequestMessage{
			Id:            c.GetTarget(),
			MessageHeader: message.MessageHeader,
			Payload:       message.Payload,
		})
	if err != nil {
		log.Println(err.Error())
		*ch <- false
		return err
	} else {
		*ch <- true
	}
	return err
}

func WaitDelay(tm int, wg *sync.WaitGroup) {
	//log.Println("Delaying send operation of", tm)
	time.Sleep(time.Duration(tm) * time.Millisecond)
	defer wg.Done()
}

//close connection
func (c *GrpcClient) Close() error {
	err := c.Client.Connection.Close()
	if err != nil {
		return err
	}
	log.Println("Connection closed")
	return nil
}

func (c *GrpcClient) GetTarget() string {
	return c.Client.Connection.Target()
}

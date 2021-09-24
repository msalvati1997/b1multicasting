package basic

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/basic/proto"
	"b1multicasting/pkg/multicasting"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Client struct {
	Client     proto.EndToEndServiceClient
	Connection *grpc.ClientConn
	delay      uint //millisecond
	address    string
}

type GrpcClient struct {
	client *Client
}

//method to connect to grpcServer
func Connect(address string, delay uint) (*GrpcClient, error) {
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
	cl.address = address
	cl.delay = delay

	return &GrpcClient{cl}, nil
}

//method to send message to GrpcServer
func (c *GrpcClient) Send(id string, message basic.Message, ch *chan bool) error {

	if c.client == nil {
		panic("Closed Connection")
	}
	//delay..
	if c.client.delay > 0 {
		delayDuration := time.Duration(c.client.delay)
		time.Sleep(delayDuration * time.Millisecond)
	}
	_, err := c.client.Client.SendMessage(context.Background(),
		&proto.RequestMessage{
			Id:            id,
			MessageHeader: message.MessageHeader,
			Payload:       message.Payload})

	if err != nil {
		log.Println(err.Error())
		*ch <- false
		return err
	} else {
		*ch <- true
	}
	return err
}

//close connection
func (c *GrpcClient) Close() error {
	err := c.client.Connection.Close()
	if err != nil {
		return err
	}
	log.Println("Connection closed")
	return nil
}

func (c *GrpcClient) GetTarget() interface{} {
	return c.client.Connection.Target()
}

func (c *GrpcClient) SendSeq(g string, m multicasting.SeqMessage, ch *chan bool) error {

	return nil
}

package utils

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"sync"
	"time"
)

var (
	GoPool Pool
)

type Pool struct {
	Mu          sync.Mutex
	MessageCh   chan basic.Message
	connections *multicasting.Conns
}

func ProcessMessage(msgChan chan basic.Message) {

	for {
		select {
		case data := <-msgChan:
			if data.MessageHeader["type"] == "TOD" {
				err := multicasting.Cnn.TODMulticast(data.MessageHeader["GroupId"], data)
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
			if data.MessageHeader["type"] == "CO" {
				err := multicasting.Cnn.COMulticast(data.MessageHeader["GroupId"], data)
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
			if data.MessageHeader["type"] == "TOC" {
				//err := multicasting.Seq.TOCMulticast(data.MessageHeader["GroupId"], data)
				err := multicasting.Seq.TOCMulticast2(data.MessageHeader["GroupId"], data)
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
		}
	}
}

func (p *Pool) Initialize(nthreads int, connections *multicasting.Conns) {

	p.connections = connections
	p.MessageCh = make(chan basic.Message, 50)

	GoPool.Mu.Lock()
	defer GoPool.Mu.Unlock()

	for i := 0; i < nthreads; i++ {
		go ProcessMessage(GoPool.MessageCh)
	}
}

package utils

import (
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	"strconv"
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
			log.Println("Processing message..")
			if data.MessageHeader["type"] == "B" {
				currentTime := time.Now()
				log.Println("Start B_SENDING")
				data.MessageHeader["ProcessId"] = utils.MystringId
				data.MessageHeader["i"] = utils.GenerateUID()
				data.MessageHeader["TIME"] = currentTime.Format("2006-01-02 15:04:05")
				err := multicasting.Cnn.BMulticast(data.MessageHeader["GroupId"], data)
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
			if data.MessageHeader["type"] == "TOD" {
				log.Println("Start TOD_SENDING")
				currentTime := time.Now()
				data.MessageHeader["TIME"] = currentTime.Format("2006-01-02 15:04:05")
				data.MessageHeader["i"] = utils.GenerateUID()
				data.MessageHeader["s"] = strconv.FormatUint(utils.Clock.Tock(), 10)
				data.MessageHeader["ProcessId"] = utils.MystringId
				err := multicasting.Cnn.TODMulticast(data.MessageHeader["GroupId"], data)
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
			if data.MessageHeader["type"] == "CO" {
				log.Println("Start CO_SENDING")
				currentTime := time.Now()
				data.MessageHeader["TIME"] = currentTime.Format("2006-01-02 15:04:05")
				data.MessageHeader["i"] = utils.GenerateUID()
				data.MessageHeader["ProcessId"] = utils.MystringId
				utils.Vectorclock.TickV(utils.Myid)
				data.MessageHeader["s"] = strconv.FormatUint(utils.Vectorclock.TockV(utils.Myid), 10)
				for p := 0; p < multicasting.GetNumbersOfClients(); p++ {
					data.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(utils.Vectorclock.TockV(p), 10)
				}
				err := multicasting.Cnn.COMulticast(data.MessageHeader["GroupId"], data)
				node := DelivererNode{NodeId: data.MessageHeader["ProcessId"]}
				Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
					Deliverer: node,
					Status:    true,
					M:         data,
				})
				log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
				if err != nil {
					go func() {
						time.Sleep(time.Second * 5)
						msgChan <- data
					}()
				}
			}
			if data.MessageHeader["type"] == "TOC" {
				currentTime := time.Now()
				data.MessageHeader["TIME"] = currentTime.Format("2006-01-02 15:04:05")
				log.Println("Start TOC_SENDING")
				data.MessageHeader["i"] = utils.GenerateUID()
				data.MessageHeader["ProcessId"] = utils.MystringId
				err := multicasting.Seq.TOCMulticast(data.MessageHeader["GroupId"], data)
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

func (p *Pool) Initialize(nthreads int) {
	log.Println("initialize pool")

	p.MessageCh = make(chan basic.Message, 50)

	GoPool.Mu.Lock()
	defer GoPool.Mu.Unlock()

	for i := 0; i < nthreads; i++ {
		go ProcessMessage(GoPool.MessageCh)
	}
}

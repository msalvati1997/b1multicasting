package multicasting

import (
	utils "b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"log"
	"strconv"
	"sync"
)

//implementing the TOD algorithm
func (c *Conns) TODMulticast(g string, m basic.Message) error {
	var wg sync.WaitGroup
	utils.Clock.Tick()
	m.MessageHeader["s"] = strconv.FormatUint(utils.Clock.Tock(), 10)
	ch := make(chan bool, len(c.conns))
	for i := 0; i < len(c.conns); i++ {
		index := i
		wg.Add(1)
		go func() {
			err := c.conns[index].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
			wg.Done()
		}()
	}
	defer wg.Wait()
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.conns[i].GetTarget())
		} else {
			//log.Println("Message correctly sent to ", c.conns[i].GetTarget())
		}
	}
	return nil
}

//Invia i messaggi di ack ai processi del relativo messaggio m
//salva in una coda i messaggi di ack relativo al messagio..può effettuare la deliver solo quando tutti hanno ricevuto il messaggio
func (c *Conns) ACKMulticast(g string, m basic.Message) error {

	ch := make(chan bool, len(c.conns))
	var wg sync.RWMutex
	for i := 0; i < len(c.conns); i++ {
		wg.Lock()
		index := i
		go func() {
			m.MessageHeader["type"] = "ACK"
			err := c.conns[index].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
			defer wg.Unlock()
		}()
	}
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.conns[i].GetTarget())
			//prova a rinviarlo
		} else {
			//log.Println("Message correctly sent to ", c.conns[i].GetTarget())
		}
	}
	return nil
}

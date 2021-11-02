package multicasting

import (
	utils "github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"log"
	"sync"
)

//implementing the TOD algorithm
func (c *Conns) TODMulticast(g string, m basic.Message) error {
	var wg sync.WaitGroup
	utils.Clock.Tick()
	ch := make(chan bool, len(c.Conns))
	wg.Add(len(c.Conns)) //parallelizing sending message
	for i := 0; i < len(c.Conns); i++ {
		index := i
		go func(w *sync.WaitGroup) {
			err := c.Conns[index].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
			w.Done()
		}(&wg)
	}
	defer wg.Wait()
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.Conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.Conns[i].GetTarget())
		} else {
			log.Println("Message", string(m.Payload), " correctly sent to ", c.Conns[i].GetTarget())
		}
	}
	return nil
}

//Send ack messages to the processes of its message m
//saves in a queue the messages of ack related to the message.. can deliver only when everyone has received the message
func (c *Conns) ACKMulticast(g string, m basic.Message) error {

	ch := make(chan bool, len(c.Conns))
	var wg sync.WaitGroup
	wg.Add(len(c.Conns))
	for i := 0; i < len(c.Conns); i++ {
		index := i
		go func(w *sync.WaitGroup) {
			err := c.Conns[index].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
			w.Done()
		}(&wg)
	}
	defer wg.Wait()
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.Conns); i++ {
		r := <-ch
		if r != true {
			log.Println("Message not arrived to nodes ", c.Conns[i].GetTarget())
		} else {
			log.Println("Message correctly sent to ", c.Conns[i].GetTarget())
		}
	}
	return nil
}

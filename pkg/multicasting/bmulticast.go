package multicasting

import (
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"log"
	"sync"
)

//IMPLEMENTING THE BMULTICAST ALGORITHM
//The basic multicast algorithm delivers messages to processes in an
//arbitrary order, due to arbitrary delays in the underlying one-to-one send operations.

//(1)To B-multicast(g, m): for each process p of the group g , send(p, m)
func (c *Conns) BMulticast(g string, m basic.Message) error {
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
			wg.Done()
		}(&wg)
	}
	defer wg.Wait()
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.Conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.Conns[i].GetTarget())
		} else {
			//log.Println("Message correctly sent to ", c.Conns[i].GetTarget())
		}
	}
	return nil
}

package multicasting

import (
	"b1multicasting/pkg/basic"
	"log"
)

//IMPLEMENTING THE BMULTICAST ALGORITHM
//The basic multicast algorithm delivers messages to processes in an
//arbitrary order, due to arbitrary delays in the underlying one-to-one send operations.

//(1)To B-multicast(g, m): for each process p of the group g , send(p, m)
func (c *Conns) BMulticast(g string, m basic.Message) error {

	ch := make(chan bool, len(c.conns))
	//da aggiungere la garanzia che i messaggi siano arrivati a destinazione tramite un pool di threads e utilizzo canali
	for i := 0; i < len(c.conns); i++ {
		i := i
		go func() {
			err := c.conns[i].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
		}()
		//inserisco la risposta nel canale gestito dalla pool
	}
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.conns[i].GetTarget())
			//prova a rinviarlo
		} else {
			log.Println("Message correctly sent to ", c.conns[i].GetTarget())
		}
	}
	return nil
}

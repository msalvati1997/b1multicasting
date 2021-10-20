package multicasting

import (
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"log"
	"sync"
)

//implementing the CO algorithm
func (c *Conns) COMulticast(g string, m basic.Message) error {
	var wg sync.WaitGroup
	ch := make(chan bool, len(c.Conns))

	wg.Add(len(c.Conns))
	for i := 0; i < len(c.Conns); i++ {
		index := i
		go func() {
			err := c.Conns[index].Send(g, m, &ch) //one to one send operation
			if err != nil {
				return
			}
			wg.Done()
		}()
	}
	defer wg.Wait()
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(c.Conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", c.Conns[i].GetTarget())
			//prova a rinviarlo
		} else {
			//log.Println("Message correctly sent to ", c.Conns[i].GetTarget())
		}
	}
	return nil
}

type COMessage struct {
	Msg    basic.Message     //Received message from i-th process
	Vector utils.VectorClock //The vector linked to the message sended from i-th process
}

func GetCOMessage(msg basic.Message, v utils.VectorClock) COMessage {
	return COMessage{
		Msg:    msg,
		Vector: v,
	}
}

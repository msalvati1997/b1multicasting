package multicasting

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"log"
	"strconv"
	"sync"
)

//implementing the CO algorithm
func (c *Conns) COMulticast(g string, m basic.Message) error {
	var wg sync.WaitGroup
	ch := make(chan bool, len(c.conns))
	utils.Vectorclock.TickV(utils.Myid) //increased the i-th clock of the vector
	for p := 0; p < len(c.conns); p++ {
		m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(utils.Vectorclock.TockV(p), 10)
	}
	for i := 0; i < len(c.conns); i++ {
		wg.Add(1)
		index := i
		go func() {
			m.MessageHeader["type"] = "CO"
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
			//prova a rinviarlo
		} else {
			//log.Println("Message correctly sent to ", c.conns[i].GetTarget())
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

package multicasting

import (
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"log"
)

//TOTAL ORDERED CENTRALIZED MULTICAST
func (seq *Sequencer) TOCMulticast2(g string, m basic.Message) error {
	//log.Println("Sending <m,i> to g : <", string(m.Payload), ",", m.MessageHeader["i"], ">")
	ch := make(chan bool, len(seq.Conns.Conns))
	for i := 0; i < len(seq.Conns.Conns); i++ {
		index := i
		go func() {
			//sends to  sequencer <m,i>
			if seq.Conns.Conns[index].GetTarget() == seq.SeqPort {
				msg := basic.Message{
					MessageHeader: make(map[string]string),
					Payload:       m.Payload,
				}
				msg.MessageHeader["member"] = "false"
				msg.MessageHeader["seq"] = "true"
				msg.MessageHeader["i"] = m.MessageHeader["i"]
				msg.MessageHeader["type"] = m.MessageHeader["type"]
				msg.MessageHeader["GroupId"] = m.MessageHeader["GroupId"]
				log.Println("Sending to seq")
				err := seq.Conns.Conns[index].Send(g, msg, &ch)
				if err != nil {
					return
				}
			} else {
				msg := basic.Message{
					MessageHeader: make(map[string]string),
					Payload:       m.Payload,
				}
				msg.MessageHeader["member"] = "true"
				msg.MessageHeader["seq"] = "false"
				msg.MessageHeader["i"] = m.MessageHeader["i"]
				msg.MessageHeader["type"] = m.MessageHeader["type"]
				msg.MessageHeader["GroupId"] = m.MessageHeader["GroupId"]
				//sends to other members
				log.Println("Sending to member")
				err := seq.Conns.Conns[index].Send(g, msg, &ch) //one to one send operation
				if err != nil {
					return
				}
			}
		}()
	}
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(seq.Conns.Conns); i++ {
		r := <-ch
		if r != true {
			log.Println("Message not arrived to nodes ", seq.Conns.Conns[i].GetTarget())
		} else {
			log.Println("Message correctly sent to ", seq.Conns.Conns[i].GetTarget())
		}
	}
	return nil
}

func (seq *Sequencer) TOCMulticast(g string, m basic.Message) error {

	ch := make(chan bool, len(seq.Conns.Conns))
	go func() {
		//invio al sequencer <m,i>
		m.MessageHeader["seq"] = "true"
		log.Println("Sending to seq")
		err := seq.SeqConn.Send(g, m, &ch)
		if err != nil {
			return
		}
	}()
	//check if the message correctly arrived to the sequencer
	r := <-ch
	if r != true {
		log.Println("Message not arrived to nodes ", seq.SeqConn.GetTarget())
	} else {
		log.Println("Message correctly sent to sequencer", seq.SeqConn.GetTarget())
	}
	return nil
}

var (
	Rg int
)

func init() {
	Rg = 0
}

func UpdateRg(S int) {
	Rg = S + 1
}

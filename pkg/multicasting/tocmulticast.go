package multicasting

import (
	"b1multicasting/pkg/basic"
	"log"
)

//TOTAL ORDERED CENTRALIZED MULTICAST
func (seq *Sequencer) TOCMulticast2(g string, m basic.Message) error {
	//log.Println("Sending <m,i> to g : <", string(m.Payload), ",", m.MessageHeader["i"], ">")
	ch := make(chan bool, len(seq.Conns.conns))
	for i := 0; i < len(seq.Conns.conns); i++ {
		index := i
		go func() {
			//invio al sequencer <m,i>
			if seq.Conns.conns[index].GetTarget() == seq.SeqPort {
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
				err := seq.Conns.conns[index].Send(g, msg, &ch)
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
				//invio agli altri membri <m,i>
				log.Println("Sending to member")
				err := seq.Conns.conns[index].Send(g, msg, &ch) //one to one send operation
				if err != nil {
					return
				}
			}
		}()
	}
	//check if the message correctly arrived to the nodes
	for i := 0; i < len(seq.Conns.conns); i++ {
		r := <-ch //lettura del canale
		if r != true {
			log.Println("Message not arrived to nodes ", seq.Conns.conns[i].GetTarget())
			//prova a rinviarlo
		} else {
			//log.Println("Message correctly sent to ", seq.Conns.conns[i].GetTarget())
		}
	}
	return nil
}

func (seq *Sequencer) TOCMulticast(g string, m basic.Message) error {

	ch := make(chan bool, len(seq.Conns.conns))
	go func() {
		//invio al sequencer <m,i>
		m.MessageHeader["seq"] = "true"
		log.Println("Sending to seq")
		err := seq.SeqConn.Send(g, m, &ch)
		if err != nil {
			return
		}
	}()
	//check if the message correctly arrived to the nodes
	r := <-ch //lettura del canale
	if r != true {
		log.Println("Message not arrived to nodes ", seq.SeqConn.GetTarget())
		//prova a rinviarlo
	} else {
		//log.Println("Message correctly sent to ", seq.Conns.conns[i].GetTarget())
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

package utils

import (
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	"math"
	"strconv"
	"sync"
)

type DelivererNode struct {
	NodeId string
}

type Delivery struct {
	Deliverer DelivererNode
	Status    bool
	M         basic.Message
}

type Deliverers struct {
	DelivererNodes []*Delivery
	M              sync.Mutex
}

var (
	Del Deliverers
)

func init() {
	Del = Deliverers{
		DelivererNodes: make([]*Delivery, 0, 100),
	}
}
func (node *DelivererNode) BDeliver(message basic.Message) {

	Del.M.Lock()
	defer Del.M.Unlock()

	Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
		Deliverer: *node,
		Status:    true,
		M:         message,
	})
	log.Println("Message correctly delivered ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
}

//Sequencer algorithm
//On B-deliver(<m,i>) : Bmulticast(g,<"order",i,sg>) , Sg=sg+1
func (node *DelivererNode) BDeliverSeq(message basic.Message) {
	log.Println("Sending from sequencer the message with order. Payload : ", string(message.Payload))
	message.MessageHeader["s"] = strconv.Itoa(multicasting.Seq.Sg)
	message.MessageHeader["order"] = "true"
	message.MessageHeader["type"] = "TOC2"
	msg := basic.NewMessage(message.MessageHeader, message.Payload)
	err := multicasting.Seq.Conns.BMulticast(multicasting.Seq.MulticastId, msg)
	multicasting.Seq.Sg = multicasting.Seq.Sg + 1 //update Sg = Sg+1
	if err != nil {
		return
	}
}

//On B-deliver(<m,i>) place <m,i> in HoldBackQueue q
func (node *DelivererNode) BDeliverMember(msg basic.Message) {
	//place msg in holdbackqueue
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		LowQ.mu.Lock()
		LowQ.LQ = append(LowQ.LQ, msg)
		defer LowQ.mu.Unlock()
		wg.Done()
	}(&wg)
	wg.Wait()
	//AppendMessageComplete(msg, msg.MessageHeader["i"], -99)
}

func (node *DelivererNode) BDeliverTOD(message basic.Message) {
	//msgi viene posto da ogni processo ricevente pj¬ in una coda locale queuej , ordinata in base al valore del timestamp.
	Del.M.Lock()
	defer Del.M.Unlock()
	id, _ := message.MessageHeader["i"]
	cl, _ := strconv.Atoi(message.MessageHeader["s"])
	utils.MuClock.Lock()
	max := math.Max(float64(utils.Clock.Tock()), float64(uint64(cl)))
	utils.Clock.Leap(uint64(max))
	utils.Clock.Tick()
	utils.MuClock.Unlock()
	AppendMessageComplete(message, id, cl)
	//Il processo ricevente invia in multicast un messaggio di Ack
	m := make(map[string]string)
	m["type"] = "ACK"
	m["i"] = message.MessageHeader["i"]
	m["s"] = strconv.Itoa(cl)
	//PrintACKQueue()
	//PrintQueue()
	//ACKMulticast
	go func() {
		err := multicasting.Cnn.ACKMulticast(node.NodeId, basic.NewMessage(m, []byte("ACK")))
		if err != nil {
			return
		}
	}()
}

func (node *DelivererNode) BDeliverCO(msg basic.Message) {
	itsvector := utils.NewVectorClock(multicasting.GetNumbersOfClients())
	for i := 0; i < multicasting.GetNumbersOfClients(); i++ {
		p := msg.MessageHeader[strconv.Itoa(i)]
		d, _ := strconv.Atoi(p)
		itsvector.LeapI(i, uint64(d))
	}
	//place <Vig,m> in holdbackqueue
	AppendMessageToCOQueue(msg, itsvector)
}

func (node *DelivererNode) TOCDeliver(msg basic.Message) {
	log.Println("TOCDELIVER of message ", string(msg.Payload), "with seq ", msg.MessageHeader["s"])
	Q.Mu.Lock()
	seq := msg.MessageHeader["s"]
	s, _ := strconv.Atoi(seq)
	Q.Q = append(Q.Q, multicasting.SeqMessage{
		Msg:  msg,
		I:    msg.MessageHeader["i"],
		Nseq: s,
	})
	Q.Mu.Unlock()
	//var wg sync.WaitGroup
	//	wg.Add(1)
	//go InsertOrder(message, message.MessageHeader["s"], &wg)
	//defer wg.Wait()
	//InsertSeq(message.MessageHeader["i"], seq, &wg)
}

//check if Vj[k]<=Vi[k] (k!=i) (second condition of deliver..)
func CheckOtherVectors(vectorclock utils.VectorClock, itsvector utils.VectorClock, i int) bool {
	b := false
	for k := 0; k < multicasting.GetNumbersOfClients(); k++ {
		if k != i {
			if itsvector.TockV(k) <= vectorclock.TockV(k) {
				b = true
			}
		}
	}
	return b
}

func TODDeliver() {
	for {
		if len(Q.Q) > 0 {

			SortingQueue()
			SortingACKQueue()
			Q.Mu.Lock()
			id := Q.Q[0].Msg.MessageHeader["i"]
			position, _ := GetMessagePosition(id)
			delnode := Q.Q[0].Msg.MessageHeader["delnode"]
			node := DelivererNode{NodeId: delnode}
			msg := Q.Q[0].Msg
			Q.Mu.Unlock()
			if (position == 0 && IsInACKQueue(id) == multicasting.GetNumbersOfClients()) && CheckOtherClocks(id) { //è in testa alla coda e tutti gli ack relativi a quel messaggio sono stati ricevuti
				Q.Mu.Lock()
				Del.M.Lock()
				ACKQueue.Mu.Lock()
				Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
					Deliverer: node,
					Status:    true,
					M:         msg,
				})
				PrintQueue()
				var wg sync.WaitGroup
				//removing from hold back queue
				go func(w *sync.WaitGroup) {
					wg.Add(1)
					for i := 0; i < len(Q.Q); i++ {
						if Q.Q[i].I == id {
							Q.Q = append(Q.Q[:i], Q.Q[i+1:]...)
						}
					}
					DeleteAckFromId(id)
					defer w.Done()
				}(&wg)
				wg.Wait()
				log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
				Q.Mu.Unlock()
				Del.M.Unlock()
				ACKQueue.Mu.Unlock()
			}
		}
	}
}

func CODeliver() {
	for {
		if len(COqueue.Q) > 0 {
			PrintCOQueue()
			msg := COqueue.Q[0].Msg
			itsvector := COqueue.Q[0].Vector
			id := msg.MessageHeader["ProcessId"]
			pid, _ := strconv.Atoi(id)
			delnode := COqueue.Q[0].Msg.MessageHeader["delnode"]
			node := DelivererNode{NodeId: delnode}
			//wait until Vj[j]=Vi[j]+1 and Vj[k]<=Vi[k] (k!=i)
			if utils.Vectorclock.TockV(pid)+1 == itsvector.TockV(pid) && CheckOtherVectors(utils.Vectorclock, itsvector, pid) {
				Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
					Deliverer: node,
					Status:    true,
					M:         msg,
				})
				log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
				//removing from hold back queue
				//COqueue = append(COqueue[:0], COqueue[1:]...)
				RemoveMessageFromQueue(msg, utils.Vectorclock)
				if pid != utils.Myid {
					utils.Vectorclock.TickV(pid)
				}
				log.Println("My vector after update")
				utils.PrintVector(utils.Vectorclock)
			} else {
				COqueue.Q = append(COqueue.Q, COqueue.Q[0])
				RemoveMessageFromQueue(msg, utils.Vectorclock)
			}
		}
	}
}

//On B-deliver(m_order=<'order',i'>) wait until <m,i> in holdbackqueue and S=rg
//TO-Deliver m , Rg=S+1
func TOCDeliver() {
	for {
		if len(Q.Q) > 0 {
			for i := 0; i < len(Q.Q); i++ {
				if Q.Q[i].Msg.MessageHeader["type"] == "TOC2" {
					seq := Q.Q[i].Nseq
					if seq == multicasting.Rg {
						Q.Mu.Lock()
						Del.M.Lock()
						id := Q.Q[i].Msg.MessageHeader["i"]
						delnode := Q.Q[i].Msg.MessageHeader["delnode"]
						node := DelivererNode{NodeId: delnode}
						msg := Q.Q[i].Msg
						Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
							Deliverer: node,
							Status:    true,
							M:         msg,
						})
						log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
						//removing from hold back queue
						var wg sync.WaitGroup
						go func(w *sync.WaitGroup) {
							wg.Add(1)
							for p := 0; p < len(Q.Q); p++ {
								if Q.Q[p].I == id {
									Q.Q = append(Q.Q[:p], Q.Q[p+1:]...)
								}
							}
							defer w.Done()
						}(&wg)
						wg.Wait()
						multicasting.UpdateRg(seq)
						Q.Mu.Unlock()
						Del.M.Unlock()
					}
				}
			}
		}
	}
}

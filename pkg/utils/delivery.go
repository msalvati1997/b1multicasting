package utils

//delivery implements deliver logic for each algorithm
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

var (
	G1 = false
	G2 = false
	G3 = false
)

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
	multicasting.Seq.Sg = multicasting.Seq.Sg + 1 //update Sg = Sg+1
	log.Println("updating seq to ", multicasting.Seq.Sg)
	go func() {
		err := multicasting.Seq.Conns.BMulticast(multicasting.Seq.MulticastId, msg)
		if err != nil {
			return
		}
	}()
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
	AppendMessageComplete(&QTOD, message, id, cl)
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
	QTOC.Mu.Lock()
	seq := msg.MessageHeader["s"]
	s, _ := strconv.Atoi(seq)
	QTOC.Q = append(QTOC.Q, multicasting.SeqMessage{
		Msg:  msg,
		I:    msg.MessageHeader["i"],
		Nseq: s,
	})
	QTOC.Mu.Unlock()
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
	log.Println("START DELIVERING THREAD")
	for {
		if len(QTOD.Q) > 0 {
			QTOD.Mu.Lock()
			SortingQueue(&QTOD)
			SortingACKQueue()
			id := QTOD.Q[0].Msg.MessageHeader["i"]
			position, _ := GetMessagePosition(&QTOD, id)
			delnode := QTOD.Q[0].Msg.MessageHeader["delnode"]
			node := DelivererNode{NodeId: delnode}
			msg := QTOD.Q[0].Msg
			QTOD.Mu.Unlock()
			if (position == 0 && IsInACKQueue(id) == multicasting.GetNumbersOfClients()) && CheckOtherClocks(id) { //è in testa alla coda e tutti gli ack relativi a quel messaggio sono stati ricevuti
				QTOD.Mu.Lock()
				Del.M.Lock()
				ACKQueue.Mu.Lock()
				Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
					Deliverer: node,
					Status:    true,
					M:         msg,
				})
				//removing from hold back queue
				QTOD.Q = append(QTOD.Q[:0], QTOD.Q[1:]...)
				DeleteAckFromId(id)
				log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
				QTOD.Mu.Unlock()
				Del.M.Unlock()
				ACKQueue.Mu.Unlock()
			}
		}
	}
}

func CODeliver() {
	log.Println("START DELIVERING THREAD")
	for {
		if len(COqueue.Q) > 0 {
			COqueue.Mu.Lock()
			msg := COqueue.Q[0].Msg
			itsvector := COqueue.Q[0].Vector
			id := msg.MessageHeader["ProcessId"]
			pid, _ := strconv.Atoi(id)
			delnode := COqueue.Q[0].Msg.MessageHeader["delnode"]
			node := DelivererNode{NodeId: delnode}
			COqueue.Mu.Unlock()
			//wait until Vj[j]=Vi[j]+1 and Vj[k]<=Vi[k] (k!=i)
			if utils.Vectorclock.TockV(pid)+1 == itsvector.TockV(pid) && CheckOtherVectors(utils.Vectorclock, itsvector, pid) {
				Del.DelivererNodes = append(Del.DelivererNodes, &Delivery{
					Deliverer: node,
					Status:    true,
					M:         msg,
				})
				log.Println("Message correctly DELIVERED ", string((Del.DelivererNodes[len(Del.DelivererNodes)-1].M).Payload))
				//removing from hold back queue
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
	log.Println("START DELIVERING THREAD")
	for {
		if len(QTOC.Q) > 0 {
			for _, value := range QTOC.Q {
				if value.Nseq == multicasting.Rg {
					Del.M.Lock()
					id := value.Msg.MessageHeader["i"]
					delnode := value.Msg.MessageHeader["delnode"]
					node := DelivererNode{NodeId: delnode}
					msg := value.Msg
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
						for p := 0; p < len(QTOC.Q); p++ {
							if QTOC.Q[p].I == id {
								QTOC.Q = append(QTOC.Q[:p], QTOC.Q[p+1:]...)
							}
						}
						defer w.Done()
					}(&wg)
					wg.Wait()
					multicasting.UpdateRg(value.Nseq)
					Del.M.Unlock()
				}
			}
		}
	}
}

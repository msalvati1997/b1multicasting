package test

import (
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"log"
	_ "math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	NMESSAGE_TOD = 5
	member_tod   = []string{"0", "1", "2"}
)

func InitData() {
	Node1.DelQueue = make([]basic.Message, 0, NMESSAGE_TOD*3)
	Node2.DelQueue = make([]basic.Message, 0, NMESSAGE_TOD*3)
	Node3.DelQueue = make([]basic.Message, 0, NMESSAGE_TOD*3)
	Node1.Clock = utils.NewClock()
	Node2.Clock = utils.NewClock()
	Node3.Clock = utils.NewClock()
	Node1.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*3)
	Node2.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*3)
	Node3.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*3)
	Node1.AckQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*NMESSAGE_TOD)
	Node2.AckQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*NMESSAGE_TOD)
	Node3.AckQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOD*NMESSAGE_TOD)
	Node1.Id = 0
	Node2.Id = 1
	Node3.Id = 2
}

//one to many tod
func Test_onetomany_tod(T1 *testing.T) {
	InitData()
	//initialize connection
	connections, err := MockConnections(member_tod, 1)
	if err != nil {
		return
	}
	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		go TODDeliver(i, &w, NMESSAGE_TOD)
	}
	defer w.Wait()
	//Multicasting messages
	for i := 0; i < NMESSAGE_TOD; i++ {
		m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
		//multicasting messages
		m.MessageHeader["type"] = "TOD"
		m.MessageHeader["ProcessId"] = "0"
		m.MessageHeader["GroupId"] = "TODMULTICAST_TEST"
		Node1.Clock.Tick()
		m.MessageHeader["s"] = strconv.FormatUint(Node1.Clock.Tock(), 10)
		err = connections.TODMulticast(m.MessageHeader["GroupId"], m)
		if err != nil {
			T1.Failed()
		}
	}
	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertTotalDelivering_tod(T1, NMESSAGE_TOD, &wg2)
	defer wg2.Wait()

	//CLEANUP
	T1.Cleanup(func() {
		multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})
}

//many to many tod
func Test_manytomany_tod(T2 *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	InitData()
	//initialize connections
	connections_1, err := MockConnections(member_tod, 1)
	if err != nil {
		T2.Failed()
	}
	connections_2, err := MockConnections(member_tod, 1)
	if err != nil {
		T2.Failed()
	}
	connections_3, err := MockConnections(member_tod, 1)
	if err != nil {
		T2.Failed()
	}
	connections := make([]*multicasting.Conns, 0, 3)
	connections = append(append(append(connections, connections_3), connections_2), connections_1)

	//Multicasting messages
	multicast := func(c *multicasting.Conns, wg *sync.WaitGroup, ProcessId string, nd int) {
		node := GetNode(nd)
		for i := 0; i < NMESSAGE_TOD; i++ {
			m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
			//multicasting messages
			m.MessageHeader["type"] = "TOD"
			m.MessageHeader["ProcessId"] = ProcessId
			m.MessageHeader["GroupId"] = "TODMULTICAST_TEST"
			node.Clock.Tick()
			m.MessageHeader["s"] = strconv.FormatUint(node.Clock.Tock(), 10)
			err = c.TODMulticast(m.MessageHeader["GroupId"], m)
			if err != nil {
				T2.Failed()
			}
		}
		defer wg.Done()
	}
	//MULTICASTING
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go multicast(connections[i], &wg, strconv.Itoa(i), i)
	}
	defer wg.Wait()

	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		log.Println("Start delivering for ", i)
		go TODDeliver(i, &w, NMESSAGE_TOD*3)
	}
	defer w.Wait()

	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertTotalDelivering_tod(T2, NMESSAGE_TOD*3, &wg2)
	defer wg2.Wait()

	//CLEANUP
	T2.Cleanup(func() {
		multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})

}

func TODDeliver(node int, w *sync.WaitGroup, nmessage int) {
	for {
		nodo := GetNode(node)
		if len(GetNode(node).DelQueue) == nmessage {
			log.Println("END OF DELIVERING  ..", node)
			w.Done()
			break
		}
		if len(nodo.HoldBackQueue.Q) > 0 {
			time.Sleep(100 * time.Millisecond)
			nodo.MuMsg.Lock()
			id := nodo.HoldBackQueue.Q[0].Msg.MessageHeader["i"]
			seq := nodo.HoldBackQueue.Q[0].Msg.MessageHeader["s"]
			s, _ := strconv.Atoi(seq)
			nodo.M = s
			var wg sync.WaitGroup
			msg := nodo.HoldBackQueue.Q[0].Msg
			p, _ := utils2.GetPositionQueue_(&nodo.HoldBackQueue, id)
			nodo.MuMsg.Unlock()
			if (p == 0 && utils2.IsInACKQueue_(&nodo.AckQueue, id) == multicasting.GetNumbersOfClients()) && CheckOtherClocks_(node, s) { //è in testa alla coda e tutti gli ack relativi a quel messaggio sono stati ricevuti E per ogni processo pk c’è un messaggio msgk in queuej con timestamp maggiore di quello di msgi
				nodo.MuMsg.Lock()
				nodo.MuDel.Lock()
				nodo.MuAck.Lock()
				log.Println("Delivering.. ", msg.MessageHeader["i"])
				utils2.PrintACKQueue_(&nodo.AckQueue, node)
				utils2.PrintQueue_(&nodo.HoldBackQueue, strconv.Itoa(node))
				nodo.DelQueue = append(nodo.DelQueue, msg)
				log.Println("Message ->", nodo.DelQueue[len(nodo.DelQueue)-1].MessageHeader["i"], " correctly delivered from ->", node)
				go func(wg *sync.WaitGroup) {
					wg.Add(1)
					nodo.HoldBackQueue.Q = append(nodo.HoldBackQueue.Q[:0], nodo.HoldBackQueue.Q[1:]...)
					for l := 0; l < len(nodo.DelQueue); l++ {
						for c := 0; c < len(nodo.AckQueue.Q); c++ {
							if nodo.DelQueue[l].MessageHeader["i"] == nodo.AckQueue.Q[c].I {
								nodo.AckQueue.Q = append(nodo.AckQueue.Q[:c], nodo.AckQueue.Q[c+1:]...)
							}
						}
					}
					defer wg.Done()
				}(&wg)
				nodo.MuAck.Unlock()
				nodo.MuMsg.Unlock()
				nodo.MuDel.Unlock()
				//Clean(node)
				time.Sleep(50 * time.Millisecond)
				//log.Println("after removing")
				//utils2.PrintACKQueue_(nodo.AckQueue, node)
				//utils2.PrintQueue_(nodo.HoldBackQueue, strconv.Itoa(node))
				wg.Wait()
			}
		}
	}
}

func IsPresent(node int, id string) bool {
	nodo := GetNode(node)
	b := false
	for i := 0; i < len(nodo.AckQueue.Q); i++ {
		if nodo.AckQueue.Q[i].I == id {
			b = true
		}
	}
	return b
}

func Clean(node int) {
	nodo := GetNode(node)
	for i := 0; i < len(nodo.DelQueue); i++ {
		for p := 0; p < len(nodo.AckQueue.Q); p++ {
			if nodo.DelQueue[i].MessageHeader["i"] == nodo.AckQueue.Q[p].I {
				nodo.AckQueue.Q = append(nodo.AckQueue.Q[:p], nodo.AckQueue.Q[p+1:]...)
			}
		}
	}
}

func CheckOtherClocks_(node int, seq int) bool {
	nodo := GetNode(node)
	if nodo == &Node1 {
		if seq <= Node2.M && seq <= Node3.M {
			return true
		}
	}
	if nodo == &Node2 {
		if seq <= Node1.M && seq <= Node3.M {
			return true
		}
	}
	if nodo == &Node3 {
		if seq <= Node2.M && seq <= Node1.M {
			return true
		}
	}
	return false
}

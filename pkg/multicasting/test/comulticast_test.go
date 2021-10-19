package test

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	client "b1multicasting/pkg/basic/client"
	"b1multicasting/pkg/multicasting"
	"fmt"
	"log"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	NMESSAGE_CO_2 = 2
	NMESSAGE_CO_1 = 10
	member_co     = []string{"0", "1", "2"}
)

func InitData_CO_1() {
	Node1_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_1)
	Node2_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_1)
	Node3_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_1)
	Node1_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_1)
	Node2_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_1)
	Node3_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_1)
	Node1_co.Id = 0
	Node2_co.Id = 1
	Node3_co.Id = 2
	Node1_co.VClock = utils.NewVectorClock(len(member_co))
	Node2_co.VClock = utils.NewVectorClock(len(member_co))
	Node3_co.VClock = utils.NewVectorClock(len(member_co))
}

func InitData_CO_2() {
	Node1_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_2*3)
	Node2_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_2*3)
	Node3_co.DelQueue = make([]basic.Message, 0, NMESSAGE_CO_2*3)
	Node1_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_2*3)
	Node2_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_2*3)
	Node3_co.HoldBackQueue.Q = make([]multicasting.COMessage, 0, NMESSAGE_CO_2*3)
	Node1_co.Id = 0
	Node2_co.Id = 1
	Node3_co.Id = 2
	Node1_co.VClock = utils.NewVectorClock(len(member_co))
	Node2_co.VClock = utils.NewVectorClock(len(member_co))
	Node3_co.VClock = utils.NewVectorClock(len(member_co))
}

func Test_onetomany_co(T1 *testing.T) {
	InitData_CO_1()
	connections, err := MockConnections(member_co, 1)
	if err != nil {
		return
	}
	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		go CODeliver(i, &w, NMESSAGE_CO_1)
	}
	defer w.Wait()
	//Multicasting messages
	for i := 0; i < NMESSAGE_CO_1; i++ {
		m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
		//multicasting messages
		m.MessageHeader["ProcessId"] = "0"
		m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
		m.MessageHeader["type"] = "CO"
		m.MessageHeader["N"] = strconv.Itoa(i)
		Node1_co.MuClock.Lock()
		m.MessageHeader["s"] = strconv.FormatUint(Node1_co.VClock.TockV(Node1_co.Id), 10)
		Node1_co.VClock.TickV(Node1_co.Id) //increased the i-th clock of the vector
		log.Println("Trying to multicast , ", m.MessageHeader)
		for p := 0; p < len(connections.Conns); p++ {
			m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node1_co.VClock.TockV(p), 10)
		}
		Node1_co.MuClock.Unlock()
		Node1_co.MuDel.Lock()
		Node1_co.DelQueue = append(Node1_co.DelQueue, m)
		log.Println("Message correctly delivered, sended from myself 1 ", m)
		Node1_co.MuDel.Unlock()
		err = connections.COMulticast(m.MessageHeader["GroupId"], m)
		if err != nil {
			T1.Failed()
		}
	}
	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertCasualDelivering_one_to_many(T1, NMESSAGE_CO_1, &wg2)
	defer wg2.Wait()

	//CLEANUP
	T1.Cleanup(func() {
		multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})
}

func Test_manytomany_co(T2 *testing.T) {
	InitData_CO_2()
	connections_1, err := MockConnections(member_co, 2)
	if err != nil {
		return
	}
	connections_2, err := MockConnections(member_co, 2)
	if err != nil {
		return
	}
	connections_3, err := MockConnections(member_co, 3)
	if err != nil {
		return
	}
	connections := make([]*multicasting.Conns, 0, 3)
	connections = append(append(append(connections, connections_3), connections_2), connections_1)
	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		go CODeliver(i, &w, NMESSAGE_CO_2*3)
	}
	defer w.Wait()
	//Multicasting messages
	multicast1 := func(wg *sync.WaitGroup) {
		for i := 0; i < len(NODE1MSG); i++ {
			m := NODE1MSG[i]
			m.MessageHeader["i"] = utils.GenerateUID()
			m.MessageHeader["ProcessId"] = "0"
			m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
			m.MessageHeader["type"] = "CO"
			Node1_co.MuClock.Lock()
			Node1_co.VClock.TickV(Node1_co.Id) //increased the i-th clock of the vector
			m.MessageHeader["s"] = strconv.FormatUint(Node1_co.VClock.TockV(Node1_co.Id), 10)
			for p := 0; p < len(connections_1.Conns); p++ {
				m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node1_co.VClock.TockV(p), 10)
			}
			log.Println("Trying to multicast , ", m.MessageHeader)
			Node1_co.MuClock.Unlock()
			Node1_co.MuDel.Lock()
			Node1_co.DelQueue = append(Node1_co.DelQueue, m)
			log.Println("Message correctly delivered, sended from myself 0 ", m)
			Node1_co.MuDel.Unlock()
			err = connections_1.COMulticast(m.MessageHeader["GroupId"], m)
			if err != nil {
				T2.Failed()
			}
		}
		log.Println("STOP MULTICASTING FOR 0")
		wg.Done()
	}

	multicast2 := func(wg *sync.WaitGroup) {
		count := 0
		for {
			if ContainsMessage(Node2_co.DelQueue, "MA") {
				count = count + 1
				m := MC
				m.MessageHeader["i"] = utils.GenerateUID()
				m.MessageHeader["ProcessId"] = "1"
				m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
				m.MessageHeader["type"] = "CO"
				Node2_co.MuClock.Lock()
				Node2_co.VClock.TickV(Node2_co.Id) //increased the i-th clock of the vector
				m.MessageHeader["s"] = strconv.FormatUint(Node2_co.VClock.TockV(Node2_co.Id), 10)
				for p := 0; p < len(connections_2.Conns); p++ {
					m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node2_co.VClock.TockV(p), 10)
				}
				log.Println("Trying to multicast , ", m.MessageHeader)
				Node2_co.MuClock.Unlock()
				Node2_co.MuDel.Lock()
				Node2_co.DelQueue = append(Node2_co.DelQueue, m)
				log.Println("Message correctly delivered, sended from myself 1 ", m)
				Node2_co.MuDel.Unlock()
				err = connections_2.COMulticast(m.MessageHeader["GroupId"], m)
				if err != nil {
					T2.Failed()
				}
				break
			}
		}
		for {
			if ContainsMessage(Node2_co.DelQueue, "MB") {
				time.Sleep(100 * time.Millisecond)
				count = count + 1
				m := MD
				m.MessageHeader["i"] = utils.GenerateUID()
				m.MessageHeader["ProcessId"] = "1"
				m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
				m.MessageHeader["type"] = "CO"
				Node2_co.MuClock.Lock()
				Node2_co.VClock.TickV(Node2_co.Id) //increased the i-th clock of the vector
				m.MessageHeader["s"] = strconv.FormatUint(Node2_co.VClock.TockV(Node2_co.Id), 10)
				for p := 0; p < len(connections_2.Conns); p++ {
					m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node2_co.VClock.TockV(p), 10)
				}
				log.Println("Trying to multicast , ", m.MessageHeader)
				Node2_co.MuClock.Unlock()
				Node2_co.MuDel.Lock()
				Node2_co.DelQueue = append(Node2_co.DelQueue, m)
				log.Println("Message correctly delivered, sended from myself 1 ", m)
				Node2_co.MuDel.Unlock()
				err = connections_2.COMulticast(m.MessageHeader["GroupId"], m)
				if err != nil {
					T2.Failed()
				}
				break
			}
		}
		if count == 2 {
			log.Println("STOP MULTICASTING FOR 1")
			wg.Done()
		}
	}

	multicast3 := func(wg *sync.WaitGroup) {
		count := 0
		for {
			if ContainsMessage(Node3_co.DelQueue, "MA") {
				m := ME
				m.MessageHeader["i"] = utils.GenerateUID()
				m.MessageHeader["ProcessId"] = "2"
				m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
				m.MessageHeader["type"] = "CO"
				Node3_co.MuClock.Lock()
				Node3_co.VClock.TickV(Node3_co.Id) //increased the i-th clock of the vector
				m.MessageHeader["s"] = strconv.FormatUint(Node3_co.VClock.TockV(Node3_co.Id), 10)
				for p := 0; p < len(connections_3.Conns); p++ {
					m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node3_co.VClock.TockV(p), 10)
				}
				Node3_co.MuClock.Unlock()
				log.Println("Trying to multicast , ", m.MessageHeader)
				Node3_co.MuDel.Lock()
				Node3_co.DelQueue = append(Node3_co.DelQueue, m)
				log.Println("Message correctly delivered, sended from myself 2 ", m)
				Node3_co.MuDel.Unlock()
				err = connections_3.COMulticast(m.MessageHeader["GroupId"], m)
				if err != nil {
					T2.Failed()
				}
				count = count + 1
				break
			}
		}
		for {
			if true {
				m := MF
				m.MessageHeader["i"] = utils.GenerateUID()
				m.MessageHeader["ProcessId"] = "2"
				m.MessageHeader["GroupId"] = "COMULTICAST_TEST"
				m.MessageHeader["type"] = "CO"
				Node3_co.MuClock.Lock()
				Node3_co.VClock.TickV(Node3_co.Id) //increased the i-th clock of the vector
				m.MessageHeader["s"] = strconv.FormatUint(Node3_co.VClock.TockV(Node3_co.Id), 10)
				for p := 0; p < len(connections_3.Conns); p++ {
					m.MessageHeader[strconv.Itoa(p)] = strconv.FormatUint(Node3_co.VClock.TockV(p), 10)
				}
				log.Println("Trying to multicast , ", m.MessageHeader)
				Node3_co.MuClock.Unlock()
				Node3_co.MuDel.Lock()
				Node3_co.DelQueue = append(Node3_co.DelQueue, m)
				log.Println("Message correctly delivered, sended from myself 2 ", m)
				Node3_co.MuDel.Unlock()
				err = connections_3.COMulticast(m.MessageHeader["GroupId"], m)
				if err != nil {
					T2.Failed()
				}
				count = count + 1
				break
			}
		}
		if count == 2 {
			log.Println("STOP MULTICASTING FOR 2")
			wg.Done()
		}
	}
	//Multicasting messages
	var wg sync.WaitGroup
	wg.Add(3)
	go multicast1(&wg)
	go multicast2(&wg)
	go multicast3(&wg)
	defer wg.Wait()

	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertCasualDelivering_Many_to_Many(T2, 6, &wg2)
	defer wg2.Wait()

	fmt.Println("The sequence of delivering of each nodes : ")
	printSequence(Node1_co.DelQueue)
	printSequence(Node2_co.DelQueue)
	printSequence(Node3_co.DelQueue)

	//CLEANUP
	T2.Cleanup(func() {
		//multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})
}

func CODeliver(i int, w *sync.WaitGroup, nmessage int) {
	log.Println("Start deliver for", i)
	for {
		nodo := GetNode_co(i)
		if len(GetNode_co(i).DelQueue) == nmessage {
			log.Println("END OF DELIVERING  ..", i)
			w.Done()
			break
		}
		if len(nodo.HoldBackQueue.Q) > 0 {
			for p := 0; p < len(nodo.HoldBackQueue.Q); p++ {
				nodo.MuMsg.Lock()
				msg := nodo.HoldBackQueue.Q[p].Msg
				itsvector := nodo.HoldBackQueue.Q[p].Vector
				id := msg.MessageHeader["ProcessId"]
				pid, _ := strconv.Atoi(id)
				nodo.MuMsg.Unlock()
				if i != pid && nodo.VClock.TockV(pid)+1 == itsvector.TockV(pid) && CheckOtherVectors(nodo.VClock, itsvector, pid) {
					nodo.MuDel.Lock()
					nodo.MuMsg.Lock()
					nodo.MuClock.Lock()
					nodo.DelQueue = append(nodo.DelQueue, msg)
					//	log.Println("My vector before update")
					//	utils.PrintVector(nodo.VClock)
					var wg sync.WaitGroup
					wg.Add(1)
					go func(w *sync.WaitGroup) {
						nodo.HoldBackQueue.Q = append(nodo.HoldBackQueue.Q[:p], nodo.HoldBackQueue.Q[p+1:]...)
						defer nodo.MuMsg.Unlock()
						defer w.Done()
					}(&wg)
					wg.Wait()
					if pid != i {
						nodo.VClock.TickV(pid)
					}
					//	log.Println("Message correctly DELIVERED from ", i, " ", string((nodo.DelQueue[len(nodo.DelQueue)-1]).Payload))
					//	log.Println("My delivering queue len ", len(nodo.DelQueue))
					//	log.Println("My queue remains  len ", len(nodo.HoldBackQueue.Q))
					//	for d:=0;d<len(nodo.HoldBackQueue.Q);d++ {
					//		utils.PrintVector(nodo.HoldBackQueue.Q[d].Vector)
					//	}
					//	log.Println("My vector after update")
					//	utils.PrintVector(nodo.VClock)
					nodo.MuClock.Unlock()
					nodo.MuDel.Unlock()
				}
			}
		}
	}
}

func CheckOtherVectors(clock utils.VectorClock, itsvector utils.VectorClock, i int) bool {
	b := false
	for k := 0; k < multicasting.GetNumbersOfClients(); k++ {
		if k != i {
			if itsvector.TockV(k) <= clock.TockV(k) {
				b = true
			} else {
				b = false
			}
		}
	}
	if b == true {
		//log.Println("Check others clock ok - receveid from ",i)
		//	utils.PrintVector(itsvector)
		//	utils.PrintVector(clock)
	}
	if b == false {
		//log.Println("Check others clock not ok - receveid from ",i)
		//utils.PrintVector(itsvector)
		//utils.PrintVector(clock)
	}
	return b
}

func ContainsMessage(queue []basic.Message, c string) bool {
	for i := 0; i < len(queue); i++ {
		if queue[i].MessageHeader["c"] == c {
			return true
		}
	}
	return false
}

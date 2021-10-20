package test

import (
	_ "github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	_ "math/rand"
	"sync"
	"testing"
)

var (
	NMESSAGE_TOC_1 = 5
	NMESSAGE_TOC_2 = 3
	member_toc     = []string{"0", "1", "2"}
)

func InitData_TOC_1() {
	Node1_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_1)
	Node2_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_1)
	Node3_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_1)
	Node1_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_1)
	Node2_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_1)
	Node3_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_1)
	Node1_toc.Rg = 0
	Node2_toc.Rg = 0
	Node3_toc.Rg = 0
	Sg = 0
	Node1_toc.SeqStruct = multicasting.Sequencer{}
	Node2_toc.SeqStruct = multicasting.Sequencer{}
	Node3_toc.SeqStruct = multicasting.Sequencer{}
}
func InitData_TOC_2() {
	Node1_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_2*3)
	Node2_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_2*3)
	Node3_toc.DelQueue = make([]basic.Message, 0, NMESSAGE_TOC_2*3)
	Node1_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_2*3)
	Node2_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_2*3)
	Node3_toc.HoldBackQueue.Q = make([]multicasting.SeqMessage, 0, NMESSAGE_TOC_2*3)
	Node1_toc.Rg = 0
	Node2_toc.Rg = 0
	Node3_toc.Rg = 0
	Sg = 0
	Node1_toc.SeqStruct = multicasting.Sequencer{}
	Node2_toc.SeqStruct = multicasting.Sequencer{}
	Node3_toc.SeqStruct = multicasting.Sequencer{}
}

//one to many toc
func Test_onetomany_toc(T1 *testing.T) {
	InitData_TOC_1()
	//initialize connection
	connections, err := MockConnections(member_toc, 2)
	if err != nil {
		return
	}
	member := make([]string, 0, 3)
	member = append(append(append(member, "0"), "1"), "2")
	sequencerPort := multicasting.SelectingSequencer(member)
	log.Println("The sequencer is ", sequencerPort)
	seqCon, err := connections.GetGrpcClient(sequencerPort)
	if err != nil {
		log.Println("Error in find connection with sequencer..", err.Error())
		return
	}

	Node2_toc.SeqStruct.SeqConn = *seqCon
	Node2_toc.SeqStruct.Conns = *connections
	Node2_toc.SeqStruct.SeqPort = sequencerPort
	Node2_toc.SeqStruct.MulticastId = "TOCMULTICAST_TEST"

	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		log.Println("Start delivering for", i)
		go TOCDeliver_(i, &w, NMESSAGE_TOC_1)
	}
	defer w.Wait()
	//Multicasting messages
	for i := 0; i < NMESSAGE_TOC_1; i++ {
		m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
		//multicasting messages
		m.MessageHeader["type"] = "TOC"
		m.MessageHeader["ProcessId"] = "0"
		m.MessageHeader["GroupId"] = "TOCMULTICAST_TEST"
		log.Println(m.MessageHeader)
		err = Node2_toc.SeqStruct.TOCMulticast(m.MessageHeader["GroupId"], m)
		if err != nil {
			T1.Failed()
		}
	}
	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertTotalDelivering_toc(T1, NMESSAGE_TOC_1, &wg2)
	defer wg2.Wait()

	//CLEANUP
	T1.Cleanup(func() {
		Node1_toc = Toc_Node{}
		Node2_toc = Toc_Node{}
		Node2_toc = Toc_Node{}
		multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})

}

//many to many tod
func Test_manytomany_toc(T2 *testing.T) {

	InitData_TOC_2()
	//initialize connections
	connections_1, err := MockConnections(member_toc, 3)
	if err != nil {
		T2.Failed()
	}
	connections_2, err := MockConnections(member_toc, 3)
	if err != nil {
		T2.Failed()
	}
	connections_3, err := MockConnections(member_toc, 3)
	if err != nil {
		T2.Failed()
	}
	sequencerPort := multicasting.SelectingSequencer(member_toc)
	log.Println("The sequencer is ", sequencerPort)
	if err != nil {
		log.Println("Error in find connection with sequencer..", err.Error())
		return
	}
	connections := make([]*multicasting.Conns, 0, 3)
	connections = append(append(append(connections, connections_3), connections_2), connections_1)
	//Multicasting messages
	multicast1 := func(wg *sync.WaitGroup) {
		multicasterId := "TOCMULTICAST_TEST"
		seqCon, _ := connections[0].GetGrpcClient(sequencerPort)
		Node1_toc.SeqStruct.Conns = *connections[0]
		Node1_toc.SeqStruct.SeqConn = *seqCon
		Node1_toc.SeqStruct.SeqPort = sequencerPort
		Node1_toc.SeqStruct.MulticastId = multicasterId
		for i := 0; i < NMESSAGE_TOC_2; i++ {
			m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
			//multicasting messages
			m.MessageHeader["type"] = "TOC"
			m.MessageHeader["ProcessId"] = "0"
			m.MessageHeader["GroupId"] = "TOCMULTICAST_TEST"
			err = Node1_toc.SeqStruct.TOCMulticast("TOCMULTICAST_TEST", m)
			if err != nil {
				T2.Failed()
			}
		}
		defer wg.Done()
	}

	multicast2 := func(wg *sync.WaitGroup) {
		multicasterId := "TOCMULTICAST_TEST"
		seqCon, _ := connections[1].GetGrpcClient(sequencerPort)
		Node2_toc.SeqStruct.Conns = *connections[1]
		Node2_toc.SeqStruct.SeqConn = *seqCon
		Node2_toc.SeqStruct.SeqPort = sequencerPort
		Node2_toc.SeqStruct.MulticastId = multicasterId
		for i := 0; i < NMESSAGE_TOC_2; i++ {
			m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
			//multicasting messages
			m.MessageHeader["type"] = "TOC"
			m.MessageHeader["ProcessId"] = "1"
			m.MessageHeader["GroupId"] = "TOCMULTICAST_TEST"
			err = Node2_toc.SeqStruct.TOCMulticast("TOCMULTICAST_TEST", m)
			if err != nil {
				T2.Failed()
			}
		}
		defer wg.Done()
	}

	multicast3 := func(wg *sync.WaitGroup) {
		multicasterId := "TOCMULTICAST_TEST"
		seqCon, _ := connections[2].GetGrpcClient(sequencerPort)
		Node3_toc.SeqStruct.Conns = *connections[2]
		Node3_toc.SeqStruct.SeqConn = *seqCon
		Node3_toc.SeqStruct.SeqPort = sequencerPort
		Node3_toc.SeqStruct.MulticastId = multicasterId
		for i := 0; i < NMESSAGE_TOC_2; i++ {
			m := basic.NewMessage(newHeader(), []byte(RandomString(16)))
			//multicasting messages
			m.MessageHeader["type"] = "TOC"
			m.MessageHeader["ProcessId"] = "2"
			m.MessageHeader["GroupId"] = "TOCMULTICAST_TEST"
			err = Node3_toc.SeqStruct.TOCMulticast("TOCMULTICAST_TEST", m)
			if err != nil {
				T2.Failed()
			}
		}
		defer wg.Done()
	}
	//Multicasting
	var wg sync.WaitGroup
	wg.Add(3)
	go multicast1(&wg)
	go multicast2(&wg)
	go multicast3(&wg)
	defer wg.Wait()
	//DELIVERING
	var w sync.WaitGroup
	w.Add(3)
	for i := 0; i < 3; i++ {
		log.Println("Start delivering for ", i)
		go TOCDeliver_(i, &w, NMESSAGE_TOC_2*3)
	}
	defer w.Wait()
	//ASSERTION
	var wg2 sync.WaitGroup
	wg2.Add(1)
	AssertTotalDelivering_toc(T2, NMESSAGE_TOC_2*3, &wg2)
	defer wg2.Wait()
	//CLEANUP
	T2.Cleanup(func() {
		Node1_toc = Toc_Node{}
		Node2_toc = Toc_Node{}
		Node2_toc = Toc_Node{}
		multicasting.Cnn.Conns = make([]*client.GrpcClient, 0, 100)
	})
}

func TOCDeliver_(node int, w *sync.WaitGroup, nmessage int) {
	for {
		nd := GetNode_toc(node)
		if len(nd.DelQueue) == nmessage {
			log.Println("END OF DELIVERING  ..", node)
			w.Done()
			break
		}
		if len(nd.HoldBackQueue.Q) > 0 {
			for i := 0; i < len(nd.HoldBackQueue.Q); i++ {
				seq := nd.HoldBackQueue.Q[i].Nseq
				if seq == nd.Rg {
					msg := nd.HoldBackQueue.Q[i].Msg
					nd.MuDel.Lock()
					nd.DelQueue = append(nd.DelQueue, msg)
					nd.MuDel.Unlock()
					log.Println("Message  correctly delivered from", node, string(nd.DelQueue[len(nd.DelQueue)-1].Payload))
					nd.MuMsg.Lock()
					nd.HoldBackQueue.Q = append(nd.HoldBackQueue.Q[:0], nd.HoldBackQueue.Q[1:]...)
					nd.MuMsg.Unlock()
					nd.Rg = nd.Rg + 1
				}
			}
		}
	}
}

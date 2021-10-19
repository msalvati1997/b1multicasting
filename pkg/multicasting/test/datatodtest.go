package test

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	utils2 "b1multicasting/pkg/utils"
	"sort"
	"strconv"
	"sync"
)

var (
	Node1 Tod_Node
	Node2 Tod_Node
	Node3 Tod_Node
)

type Tod_Node struct {
	Id            int
	HoldBackQueue utils2.Queue
	AckQueue      utils2.ACKQ
	Clock         utils.LogicalClock
	DelQueue      []basic.Message
	M             int
	MuMsg         sync.Mutex
	MuAck         sync.Mutex
	MuDel         sync.Mutex
	MuClock       sync.Mutex
}

func AppendMessageTOD(message basic.Message, id string, cl uint64, i int) {
	nodo := GetNode(i)
	message.MessageHeader["deliverable"] = "false"
	nodo.HoldBackQueue.Q = append(nodo.HoldBackQueue.Q, multicasting.SeqMessage{
		Msg:  message,
		I:    id,
		Nseq: int(cl),
	})
	sort.SliceStable(nodo.HoldBackQueue.Q, func(i, j int) bool {
		if nodo.HoldBackQueue.Q[i].Nseq == nodo.HoldBackQueue.Q[j].Nseq {
			n1, _ := strconv.Atoi(nodo.HoldBackQueue.Q[i].Msg.MessageHeader["ProcessId"])
			n2, _ := strconv.Atoi(nodo.HoldBackQueue.Q[j].Msg.MessageHeader["ProcessId"])
			return (nodo.HoldBackQueue.Q[i].Nseq <= nodo.HoldBackQueue.Q[j].Nseq) && n1 < n2
		}
		return nodo.HoldBackQueue.Q[i].Nseq < nodo.HoldBackQueue.Q[j].Nseq
	})
	utils2.PrintQueue_(&nodo.HoldBackQueue, strconv.Itoa(i))
	nodo.MuMsg.Unlock()
}

func AppendACK(msg basic.Message, index string, nseq int, node int) {
	nodo := GetNode(node)
	nodo.MuAck.Lock()
	nodo.AckQueue.Q = append(nodo.AckQueue.Q, multicasting.SeqMessage{
		Msg:  msg,
		I:    index,
		Nseq: nseq,
	})
	nodo.MuAck.Unlock()
}

func GetNode(id int) *Tod_Node {
	if id == 0 {
		return &Node1
	}
	if id == 1 {
		return &Node2
	}
	if id == 2 {
		return &Node3
	}
	return nil
}

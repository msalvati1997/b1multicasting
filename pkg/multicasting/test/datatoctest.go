package test

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	utils2 "b1multicasting/pkg/utils"
	"sort"
	"sync"
)

var (
	Node1_toc Toc_Node
	Node2_toc Toc_Node
	Node3_toc Toc_Node
)

var (
	Sg = 0
	M  sync.Mutex
)

type Toc_Node struct {
	Id            int
	HoldBackQueue utils2.Queue
	DelQueue      []basic.Message
	SeqStruct     multicasting.Sequencer
	Rg            int
	MuMsg         sync.Mutex
	MuAck         sync.Mutex
	MuDel         sync.Mutex
}

func AppendMessageTct(message basic.Message, id string, cl uint64, i int) {
	node := GetNode_toc(i)
	node.MuMsg.Lock()
	node.HoldBackQueue.Q = append(node.HoldBackQueue.Q, multicasting.SeqMessage{
		Msg:  message,
		I:    id,
		Nseq: int(cl),
	})
	sort.Slice(node.HoldBackQueue.Q, func(p, j int) bool {
		return node.HoldBackQueue.Q[p].Nseq < node.HoldBackQueue.Q[j].Nseq
	})
	node.MuMsg.Unlock()
}

func GetNode_toc(id int) *Toc_Node {
	if id == 0 {
		return &Node1_toc
	}
	if id == 1 {
		return &Node2_toc
	}
	if id == 2 {
		return &Node3_toc
	}
	return nil
}

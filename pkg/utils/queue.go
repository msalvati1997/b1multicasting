package utils

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"errors"
	"log"
	"sort"
	"strconv"
	"sync"
)

type Queue struct {
	Q []multicasting.SeqMessage
}
type LowPriorityQueue struct {
	LQ []basic.Message
	mu sync.RWMutex
}

var (
	LowQ LowPriorityQueue
)
var (
	Q Queue
)

func init() {
	Q.Q = make([]multicasting.SeqMessage, 0, 100)
	LowQ.LQ = make([]basic.Message, 100)
}

func SortingQueue() {

	sort.Slice(Q.Q, func(i, j int) bool {
		return Q.Q[i].Nseq < Q.Q[j].Nseq
	})
}

func InsertOrder(msg basic.Message, seq string, wg *sync.WaitGroup) {
	seq_int, _ := strconv.Atoi(seq)
	go func() {
		for i := 0; i < len(LowQ.LQ); i++ {
			if LowQ.LQ[i].MessageHeader["i"] == msg.MessageHeader["i"] {
				AppendMessageComplete(LowQ.LQ[i], LowQ.LQ[i].MessageHeader["i"], seq_int)
			}
		}
	}()
	wg.Done()
}

func PrintQueue() {

	for i := 0; i < len(Q.Q); i++ {
		log.Println("<", string(Q.Q[i].Msg.Payload), ",", Q.Q[i].I, ",", Q.Q[i].Nseq, ">")
	}
}

func AppendMessageComplete(message basic.Message, i string, seq int) {

	seqm := multicasting.SeqMessage{Msg: message, I: i, Nseq: seq}
	Q.Q = append(Q.Q, seqm)
	SortingQueue()
}

func GetMinClockToDeliver() int {
	SortingQueue()
	return Q.Q[0].Nseq
}

func InsertSeq(id string, nseq int, wg *sync.WaitGroup) bool {
	defer wg.Done()

	for i := 0; i < len(Q.Q); i++ {
		if Q.Q[i].I == id {
			(&Q.Q[i]).Nseq = nseq
			return true
		}
	}
	return false
}

func DeleteMessage(id string, wg *sync.WaitGroup) error {
	pos := 0
	for i := 0; i < len(Q.Q); i++ {
		if Q.Q[i].I == id {
			pos = i
		}
	}
	Q.Q = append(Q.Q[:pos], Q.Q[pos+1:]...)
	defer wg.Done()
	return nil
}

func GetMessagePosition(id string) (int, error) {

	for i := 0; i < len(Q.Q); i++ {
		if Q.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The message is not in the queue")
}

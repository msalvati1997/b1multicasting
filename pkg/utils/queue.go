package utils

import (
	"errors"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	"sort"
	"strconv"
	"sync"
)

type Queue struct {
	Q  []multicasting.SeqMessage
	Mu sync.Mutex
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

func PrintQueue() {
	for i := 0; i < len(Q.Q); i++ {
		log.Println("<", string(Q.Q[i].Msg.Payload), ",", Q.Q[i].I, ",", Q.Q[i].Nseq, ">")
	}
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

func PrintQueue_(queue *Queue, id string) {
	log.Println("Print queue of ", id)
	for i := 0; i < len(queue.Q); i++ {
		log.Println("<", string(queue.Q[i].Msg.Payload), ",", queue.Q[i].I, ",", queue.Q[i].Nseq, ">")
	}
}

func AppendMessageComplete(message basic.Message, i string, seq int) {
	seqm := multicasting.SeqMessage{Msg: message, I: i, Nseq: seq}
	Q.Q = append(Q.Q, seqm)
	sort.Slice(Q.Q, func(i, j int) bool {
		return Q.Q[i].Nseq < Q.Q[j].Nseq
	})
}

func AppendMessageComplete_(queue *Queue, message basic.Message, i string, seq int) {
	seqm := multicasting.SeqMessage{Msg: message, I: i, Nseq: seq}
	queue.Q = append(queue.Q, seqm)
	SortingQueue_(queue)
}

func SortingQueue_(queue *Queue) {
	sort.Slice(queue.Q, func(i, j int) bool {
		return queue.Q[i].Nseq < queue.Q[j].Nseq
	})
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

func DeleteMessage_(queue *Queue, id string, wg *sync.WaitGroup) error {
	pos := 0
	for i := 0; i < len(queue.Q); i++ {
		if queue.Q[i].I == id {
			pos = i
		}
	}
	queue.Q = append(queue.Q[:pos], queue.Q[pos+1:]...)
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

func GetMessagePosition_(queue *Queue, id string) (int, error) {
	for i := 0; i < len(queue.Q); i++ {
		if queue.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The message is not in the queue")
}

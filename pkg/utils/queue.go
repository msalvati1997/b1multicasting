package utils

import (
	"errors"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/multicasting"
	"log"
	"sort"
	"sync"
)

//local queue that keeps received messages with their timestampts
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
	QTOC Queue
	QTOD Queue
)

func init() {
	QTOC.Q = make([]multicasting.SeqMessage, 0, 100)
	QTOD.Q = make([]multicasting.SeqMessage, 0, 100)
	LowQ.LQ = make([]basic.Message, 100)
}

func SortingQueue(Q *Queue) {
	sort.Slice(Q.Q, func(i, j int) bool {
		return Q.Q[i].Nseq < Q.Q[j].Nseq
	})
}

func PrintQueue(Q *Queue) {
	for i := 0; i < len(Q.Q); i++ {
		log.Println("<", string(Q.Q[i].Msg.Payload), ",", Q.Q[i].I, ",", Q.Q[i].Nseq, ">")
	}
}

func PrintQueue_(queue *Queue, id string) {
	log.Println("Print queue of ", id)
	for i := 0; i < len(queue.Q); i++ {
		log.Println("<", string(queue.Q[i].Msg.Payload), ",", queue.Q[i].I, ",", queue.Q[i].Nseq, ">")
	}
}

func AppendMessageComplete(Q *Queue, message basic.Message, i string, seq int) {
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

func InsertSeq(Q *Queue, id string, nseq int, wg *sync.WaitGroup) bool {
	defer wg.Done()
	for i := 0; i < len(Q.Q); i++ {
		if Q.Q[i].I == id {
			(&Q.Q[i]).Nseq = nseq
			return true
		}
	}
	return false
}

func DeleteMessage(Q *Queue, id string, wg *sync.WaitGroup) error {
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

func DeleteMessage_(queue *Queue, id string) error {
	queue.Mu.Lock()
	defer queue.Mu.Unlock()
	pos := 0
	for i := 0; i < len(queue.Q); i++ {
		if queue.Q[i].I == id {
			pos = i
		}
	}
	queue.Q = append(queue.Q[:pos], queue.Q[pos+1:]...)
	return nil
}

func GetMessagePosition(Q *Queue, id string) (int, error) {

	for i := 0; i < len(Q.Q); i++ {
		if Q.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The message is not in the queue")
}

func GetMessagePosition_(Q *Queue, queue *Queue, id string) (int, error) {
	for i := 0; i < len(queue.Q); i++ {
		if queue.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The message is not in the queue")
}

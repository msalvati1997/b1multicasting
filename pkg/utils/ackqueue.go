package utils

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"errors"
	"log"
	"sort"
	"sync"
)

var (
	ACKQueue ACKQ
)

type ACKQ struct {
	Q  []multicasting.SeqMessage
	Mu sync.Mutex
}

func init() {
	ACKQueue.Q = make([]multicasting.SeqMessage, 0, 100)
}

func IsInACKQueue(id string) int {
	p := 0
	for i := 0; i < len(ACKQueue.Q); i++ {
		if ACKQueue.Q[i].I == id {
			p = p + 1
		}
	}
	//log.Println("Number of ack : ", p)
	//PrintACKQueue()
	return p
}

func IsInACKQueue_(ac *ACKQ, id string) int {
	p := 0
	for i := 0; i < len(ac.Q); i++ {
		if ac.Q[i].I == id {
			p = p + 1
		}
	}
	//if p==3 {
	//	log.Println("Number of ack : ", p)
	//	}
	//	PrintACKQueue()
	return p
}

func GetPosition(id string) (int, error) {
	for i := 0; i < len(ACKQueue.Q); i++ {
		if ACKQueue.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The element is not in the ACKQueue")
}

func IsInACKQueueB(id string) bool {
	for i := 0; i < len(ACKQueue.Q); i++ {
		if ACKQueue.Q[i].I == id {
			return true
		}
	}
	return false
}

func IsInACKQueueB_(acq *ACKQ, id string) bool {
	for i := 0; i < len(acq.Q); i++ {
		if acq.Q[i].I == id {
			return true
		}
	}
	return false
}

func AppendACK(message basic.Message, i string, nseq int) {
	seqm := multicasting.SeqMessage{Msg: message, I: i, Nseq: nseq}
	ACKQueue.Q = append(ACKQueue.Q, seqm)
}

func PrintACKQueue() {
	for i := 0; i < len(ACKQueue.Q); i++ {
		log.Println("<", string(ACKQueue.Q[i].Msg.Payload), ",", ACKQueue.Q[i].I, ",", ACKQueue.Q[i].Nseq, ">")
	}
}

func DeleteAckFromId(id string) bool {
	for IsInACKQueueB(id) {
		pos, err := GetPosition(id)
		if err != nil {
			return false
		} else {
			ACKQueue.Q = append(ACKQueue.Q[:pos], ACKQueue.Q[pos+1:]...)
		}
	}
	return true
}

func DeleteAckFromId_(A *ACKQ, id string, wg *sync.WaitGroup) bool {
	defer wg.Wait()
	for IsInACKQueueB_(A, id) {
		p, err := GetPositionAck_(A, id)
		if err != nil {
			return false
		} else {
			A.Q = append(A.Q[:p], A.Q[p+1:]...)
		}
	}
	return true
}

//controlla dagli ack ricevuti se gli altri processi non abbiano un valore di timestamp minore di quello relativo al messaggio con id
func CheckOtherClocks(id string) bool {
	SortingACKQueue()
	i, err := GetPosition(id)
	if err != nil {
		return false
	}
	if ACKQueue.Q[i].Nseq <= ACKQueue.Q[0].Nseq {
		return true
	}
	return false
}

func GetPositionAck_(AC *ACKQ, id string) (int, error) {
	for i := 0; i < len(AC.Q); i++ {
		index := i
		if AC.Q[i].I == id {
			return index, nil
		}
	}
	return -1, errors.New("The element is not in the ACKQueue")
}

func GetPositionQueue_(queue *Queue, id string) (int, error) {
	for i := 0; i < len(queue.Q); i++ {
		if queue.Q[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The element is not in the Queue")
}

func SortingACKQueue() {

	sort.SliceStable(ACKQueue.Q, func(i, j int) bool {
		return ACKQueue.Q[i].Nseq < ACKQueue.Q[j].Nseq
	})
}
func PrintACKQueue_(ac *ACKQ, i int) {
	log.Println("Print ACK queue of ", i)
	for i := 0; i < len(ac.Q); i++ {
		index := i
		log.Println("<", string(ac.Q[index].Msg.Payload), ",", ac.Q[index].I, ",", ac.Q[index].Nseq, ">")
	}
}

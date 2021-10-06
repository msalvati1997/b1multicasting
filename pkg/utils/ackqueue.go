package utils

import (
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"errors"
	"log"
	"sort"
)

var (
	ACKQueue []multicasting.SeqMessage
)

func init() {
	ACKQueue = make([]multicasting.SeqMessage, 0, 100)
}

func IsInACKQueue(id string) int {
	p := 0
	for i := 0; i < len(ACKQueue); i++ {
		if ACKQueue[i].I == id {
			p = p + 1
		}
	}
	//log.Println("Number of ack : ", p)
	//	PrintACKQueue()
	return p
}

func GetPosition(id string) (int, error) {
	for i := 0; i < len(ACKQueue); i++ {
		if ACKQueue[i].I == id {
			return i, nil
		}
	}
	return -1, errors.New("The element is not in the ACKQueue")
}

func IsInACKQueueB(id string) bool {
	for i := 0; i < len(ACKQueue); i++ {
		if ACKQueue[i].I == id {
			return true
		}
	}
	return false
}

func AppendACK(message basic.Message, i string, nseq int) {
	seqm := multicasting.SeqMessage{Msg: message, I: i, Nseq: nseq}
	ACKQueue = append(ACKQueue, seqm)
}

func PrintACKQueue() {
	for i := 0; i < len(ACKQueue); i++ {
		log.Println("<", string(ACKQueue[i].Msg.Payload), ",", ACKQueue[i].I, ",", ACKQueue[i].Nseq, ">")
	}
}

func DeleteAckFromId(id string) bool {
	for IsInACKQueueB(id) {
		pos, err := GetPosition(id)
		if err != nil {
			return false
		} else {
			ACKQueue = append(ACKQueue[:pos], ACKQueue[pos+1:]...)
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
	if ACKQueue[i].Nseq <= ACKQueue[0].Nseq {
		return true
	}
	return false
}

func SortingACKQueue() {
	sort.Slice(ACKQueue, func(i, j int) bool {
		return ACKQueue[i].Nseq < ACKQueue[j].Nseq
	})
}

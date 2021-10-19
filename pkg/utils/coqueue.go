package utils

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"errors"
	"log"
	"sync"
)

var (
	COqueue COQueue
)

type COQueue struct {
	Q  []multicasting.COMessage
	Mu sync.Mutex
}

func init() {
	COqueue.Q = make([]multicasting.COMessage, 0, 100)
}

func AppendMessageToCOQueue(msg basic.Message, v utils.VectorClock) {
	COqueue.Mu.Lock()
	defer COqueue.Mu.Unlock()
	COqueue.Q = append(COqueue.Q, multicasting.COMessage{
		Msg:    msg,
		Vector: v,
	})
}

func GetVectorPosition(msg basic.Message, v utils.VectorClock) (int, error) {
	Cmsg := multicasting.GetCOMessage(msg, v)
	for i := 0; i < len(COqueue.Q); i++ {
		if COqueue.Q[i].Msg.MessageHeader["i"] == Cmsg.Msg.MessageHeader["i"] {
			return i, nil
		}
	}
	return -1, errors.New("The message m was not found on the Vector")
}

func GetVectorPosition_(CQ *COQueue, msg basic.Message) (int, error) {
	for i := 0; i < len(CQ.Q); i++ {
		if CQ.Q[i].Msg.MessageHeader["i"] == msg.MessageHeader["i"] {
			return i, nil
		}
	}
	return -1, errors.New("The message m was not found on the Queue")
}

func RemoveMessageFromQueue_(CO *COQueue, msg basic.Message, v utils.VectorClock) {
	CO.Mu.Lock()
	defer CO.Mu.Unlock()
	p, err := GetVectorPosition_(CO, msg)
	if err != nil {
		return
	}
	CO.Q = append(CO.Q[:p], CO.Q[p+1:]...)
}

func IsInQueue(co *COQueue, msg basic.Message) bool {
	for i := 0; i < len(co.Q); i++ {
		if co.Q[i].Msg.MessageHeader["i"] == msg.MessageHeader["i"] {
			return true
		}
	}
	return false
}

func RemoveMessageFromQueue(msg basic.Message, v utils.VectorClock) {
	COqueue.Mu.Lock()
	defer COqueue.Mu.Unlock()
	p, err := GetVectorPosition(msg, v)
	if err != nil {
		return
	}
	COqueue.Q = append(COqueue.Q[:p], COqueue.Q[p+1:]...)
}

func PrintCOQueue() {
	for i := 0; i < len(COqueue.Q); i++ {
		log.Println(COqueue.Q[i].Msg.MessageHeader)
		utils.PrintVector(COqueue.Q[i].Vector)
	}
}

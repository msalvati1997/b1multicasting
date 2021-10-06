package utils

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"b1multicasting/pkg/multicasting"
	"errors"
	"log"
)

var (
	COqueue []multicasting.COMessage
)

func init() {
	COqueue = make([]multicasting.COMessage, 0, 100)
}

func AppendMessageToCOQueue(msg basic.Message, v utils.VectorClock) {
	COqueue = append(COqueue, multicasting.COMessage{
		Msg:    msg,
		Vector: v,
	})
}

func GetVectorPosition(msg basic.Message, v utils.VectorClock) (int, error) {
	Cmsg := multicasting.GetCOMessage(msg, v)
	for i := 0; i < len(COqueue); i++ {
		if COqueue[i].Msg.MessageHeader["i"] == Cmsg.Msg.MessageHeader["i"] {
			return i, nil
		}
	}
	return -1, errors.New("The message m was not found on the Vector")
}

func RemoveMessageFromQueue(msg basic.Message, v utils.VectorClock) {
	p, err := GetVectorPosition(msg, v)
	if err != nil {
		return
	}
	COqueue = append(COqueue[:p], COqueue[p+1:]...)
}

func PrintCOQueue() {
	for i := 0; i < len(COqueue); i++ {
		log.Println(COqueue[i].Msg.MessageHeader)
		utils.PrintVector(COqueue[i].Vector)
	}
}

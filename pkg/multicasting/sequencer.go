package multicasting

import (
	"b1multicasting/pkg/basic"
	client "b1multicasting/pkg/basic/client"
	"math/rand"
	"sort"
)

// Sequencer
type Sequencer struct {
	SeqPort     string
	SeqConn     client.GrpcClient //sequencer connection
	MulticastId string
	Conns       Conns //groups connection
	Sg          int
	B           bool
}

type SeqMessage struct {
	Msg  basic.Message
	I    string
	Nseq int
}

var (
	Seq Sequencer
)

func init() {
	Seq.Sg = 0
}

func SelectingSequencer(member []string) string {
	rand.Seed(int64(len(member)))
	index := rand.Intn(len(member))
	sort.Strings(member)
	return member[index]
}

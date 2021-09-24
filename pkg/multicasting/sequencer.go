package multicasting

import (
	base "b1multicasting/pkg/basic"
)

// Sequencer
type Sequencer struct {
}

type SeqMessage struct {
	Msg    base.Message
	Seqnum int
}

var (
	sm SeqMessage
)

func init() {
	sm.Seqnum = 0 //initializzo il numero di sequenza a 0
}

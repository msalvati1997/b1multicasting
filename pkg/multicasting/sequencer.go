package multicasting

import (
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	"math/rand"
	"sort"
	"strconv"
	"strings"
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

func SelectingSequencer(member []string, b bool) string {
	index := 0
	sort.Strings(member)
	members_ := make(map[string]string)
	targets := GetTargets()
	n := []int{}
	if b {
		for _, k := range targets {
			splitted_strings := strings.Split(strings.Split(k, ":")[0], ".")
			in := splitted_strings[len(splitted_strings)-1]
			ind, _ := strconv.Atoi(in)
			members_[in] = k
			n = append(n, ind)
		}
		sort.Ints(n)
		rand.Seed(int64(len(Cnn.Conns)))
		index = rand.Intn(len(Cnn.Conns))
		seq := n[index]
		return members_[strconv.Itoa(seq)]
	} else {
		rand.Seed(int64(len(Cnn.Conns)))
		index = rand.Intn(len(Cnn.Conns))
		sort.Strings(member)
		return member[index]
	}
}

func GetTargets() []string {
	targets := []string{}
	for i := 0; i < len(Cnn.Conns); i++ {
		targets = append(targets, Cnn.Conns[i].GetTarget())
	}
	return targets
}

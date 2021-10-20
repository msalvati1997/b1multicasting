package test

import (
	"github.com/msalvati1997/b1multicasting/internal/utils"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	utils2 "github.com/msalvati1997/b1multicasting/pkg/utils"
	"sync"
)

var (
	Node1_co Co_Node
	Node2_co Co_Node
	Node3_co Co_Node
)

type Co_Node struct {
	Id            int
	HoldBackQueue utils2.COQueue
	DelQueue      []basic.Message
	VClock        utils.VectorClock
	MuMsg         sync.Mutex
	MuDel         sync.Mutex
	MuClock       sync.Mutex
}

func GetNode_co(id int) *Co_Node {
	if id == 0 {
		return &Node1_co
	}
	if id == 1 {
		return &Node2_co
	}
	if id == 2 {
		return &Node3_co
	}
	return nil
}

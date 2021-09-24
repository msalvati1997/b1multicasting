package utils

import (
	"b1multicasting/pkg/basic"
	"log"
	"sync"
)

type DelivererNode struct {
	NodeId string
}

type Delivery struct {
	Deliverer DelivererNode
	Status    bool //lo stato della consegna passa a true quando è stato consegnato
	M         basic.Message
}

type Deliverers struct {
	DelivererNodes []*Delivery
	M              sync.Mutex
}

var (
	del Deliverers
)

func init() {
	del = Deliverers{
		DelivererNodes: make([]*Delivery, 100),
	}
}
func (node *DelivererNode) BDeliver(message basic.Message) {

	del.M.Lock()
	defer del.M.Unlock()

	del.DelivererNodes = append(del.DelivererNodes, &Delivery{
		Deliverer: *node,
		Status:    true,
		M:         message,
	})
	log.Println("Message correctly delivered ", string((del.DelivererNodes[len(del.DelivererNodes)-1].M).Payload))
}

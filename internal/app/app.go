package app

import (
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	"github.com/msalvati1997/b1multicasting/pkg/registry/proto"
	"sync"
)

var (
	RegistryClient  proto.RegistryClient
	GMu             sync.RWMutex
	MulticastGroups map[string]*MulticastGroup
	GrpcPort        uint
)

func init() {
	MulticastGroups = make(map[string]*MulticastGroup)
}

// MulticastGroup manages the metadata associated with a group in which the node is registered
type MulticastGroup struct {
	ClientId  string
	Group     *MulticastInfo
	GroupMu   sync.RWMutex
	Messages  []basic.Message
	MessageMu sync.RWMutex
}

type MulticastInfo struct {
	MulticastId      string            `json:"multicast_id"`
	MulticastType    string            `json:"multicast_type"`
	ReceivedMessages int               `json:"received_messages"`
	Status           string            `json:"status"`
	Members          map[string]Member `json:"members"`
	Error            error             `json:"error"`
}

type Member struct {
	MemberId string `json:"member_id"`
	Address  string `json:"address"`
	Ready    bool   `json:"ready"`
}

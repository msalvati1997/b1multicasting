package reg

import (
	"github.com/msalvati1997/b1multicasting/pkg/reg/proto"
)

var MulticastType = map[string]proto.MulticastType{
	"BMULTICAST":   proto.MulticastType_BMULTICAST,
	"TOCMULTICAST": proto.MulticastType_TOCMULTICAST,
	"TODMULTICAST": proto.MulticastType_TODMULTICAST,
	"COMULTICAST":  proto.MulticastType_COMULTICAST,
}

var (
	MulticastTypes []string
)

func init() {
	MulticastTypes = make([]string, 0)

	for multicastType := range MulticastType {
		MulticastTypes = append(MulticastTypes, multicastType)
	}
}

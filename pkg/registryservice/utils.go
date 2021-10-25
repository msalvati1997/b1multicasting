package registryservice

import (
	"github.com/msalvati1997/b1multicasting/pkg/registryservice/protoregistry"
)

var MulticastType = map[string]protoregistry.MulticastType{
	"BMULTICAST":   protoregistry.MulticastType_BMULTICAST,
	"TOCMULTICAST": protoregistry.MulticastType_TOCMULTICAST,
	"TODMULTICAST": protoregistry.MulticastType_TODMULTICAST,
	"COMULTICAST":  protoregistry.MulticastType_COMULTICAST,
}

var (
	MulticastTypes []string
)

func init() {
	MulticastTypes = make([]string, 4)

	for multicastType := range MulticastType {
		MulticastTypes = append(MulticastTypes, multicastType)
	}
}

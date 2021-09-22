package multicasting

import (
	"b1multicasting/pkg/basic"
)

//Multicast API
//X IS ONE OF - B(BASIC)/CO(CASUAL)/TO(TOTAL)

type Multicast interface {
	XMulticast(g string, m basic.Message)
	XDeliver(message basic.Message)
}

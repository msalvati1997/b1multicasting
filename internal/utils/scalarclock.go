package utils

import (
	"sync/atomic"
)

var (
	Clock LogicalClock
)

func init() {
	Clock = NewClock()
}

// LogicalClock provide the timestamp for a single peer.
type LogicalClock interface {
	// Tick the clock is increased.
	Tick()

	// Tock the value present on the clock is retrieved.
	Tock() uint64

	// Leap the value on the clock leaps to the given value.
	Leap(to uint64)
}

// ProcessClock for a single process, implements the LogicalClock interface.
type ProcessClock struct {
	// Logical operation index.
	Index uint64
}

// Tick Implements the LogicalClock interface.
func (p *ProcessClock) Tick() {
	atomic.AddUint64(&p.Index, 1)
}

// Tock Implements the LogicalClock interface.
func (p *ProcessClock) Tock() uint64 {
	return atomic.LoadUint64(&p.Index)
}

// Leap Implements the LogicalClock interface.
func (p *ProcessClock) Leap(to uint64) {
	atomic.StoreUint64(&p.Index, to)
}
func NewClock() LogicalClock {
	return &ProcessClock{
		Index: 0,
	}
}

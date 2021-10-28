package utils

import (
	"log"
	"strconv"
	"sync"
)

type VectorClock interface {
	// Tick the i-th clock in the vector is increased.
	TickV(i int)
	//Tock the value present on the i-th clock is retrieved.
	TockV(i int) uint64
	//Get Lenght of the VectorClock
	GetLenght() int
	// Leap the value on the ith-clock leaps to the given value.
	LeapI(i int, to uint64)
}

var (
	Vectorclock   VectorClock
	MuVectorClock sync.Mutex
)

var (
	Myid          int
	MyAdress      string
	MyAdressSplit string
	Mymu          sync.Mutex
)

// ProcessVectorClock for a single process, implements the VectorClock interface.
type ProcessVectorClock struct {
	Vector map[int]LogicalClock
	n      int
}

func (p ProcessVectorClock) LeapI(i int, to uint64) {
	p.Vector[i].Leap(to)
}

func (p ProcessVectorClock) GetLenght() int {
	return p.n
}

func (p ProcessVectorClock) TickV(i int) {
	p.Vector[i].Tick()
}

func (p ProcessVectorClock) TockV(i int) uint64 {
	return p.Vector[i].Tock()
}

//initialization of vector clock
func NewVectorClock(n int) VectorClock {
	v := make(map[int]LogicalClock, n)
	for i := 0; i < n; i++ {
		v[i] = NewClock()
	}
	return &ProcessVectorClock{
		Vector: v,
		n:      n,
	}
}

func PrintVector(v VectorClock) {
	s := "[ "
	for i := 0; i < v.GetLenght(); i++ {
		s = s + strconv.FormatUint(v.TockV(i), 10)
		s = s + " "
	}
	s = s + "]"
	log.Println(s)
}

package test

import (
	"b1multicasting/internal/utils"
	"b1multicasting/pkg/basic"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"sync"
	"testing"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func newHeader() map[string]string {
	m := make(map[string]string)
	m["i"] = utils.GenerateUID()
	return m
}

// AssertTotalDelivering_tod
//  Testing TOD
func AssertTotalDelivering_tod(T *testing.T, nm int, wg *sync.WaitGroup) {
	for {
		del1 := Node1.DelQueue
		del2 := Node2.DelQueue
		del3 := Node3.DelQueue
		if len(del1) == nm && len(del2) == nm && len(del3) == nm {
			log.Println("The test assertion..")
			for i := 0; i < nm; i++ {
				id1 := del1[i].MessageHeader["i"]
				id2 := del2[i].MessageHeader["i"]
				id3 := del3[i].MessageHeader["i"]
				assert.Equal(T, id1, id2, "Error in comparing delivering ")
				assert.Equal(T, id2, id3, "Error in comparing delivering ")
				assert.Equal(T, id1, id3, "Error in comparing delivering ")
			}
			wg.Done()
			break
		}
	}
}

// AssertTotalDelivering_toc
//  Testing TOC
func AssertTotalDelivering_toc(T *testing.T, nm int, wg *sync.WaitGroup) {
	for {
		del1 := Node1_toc.DelQueue
		del2 := Node2_toc.DelQueue
		del3 := Node3_toc.DelQueue
		if len(del1) == nm && len(del2) == nm && len(del3) == nm {
			log.Println("The test assertion..")
			for i := 0; i < nm; i++ {
				id1 := del1[i].MessageHeader["i"]
				id2 := del2[i].MessageHeader["i"]
				id3 := del3[i].MessageHeader["i"]
				assert.Equal(T, id1, id2, "Error in comparing delivering ")
				assert.Equal(T, id2, id3, "Error in comparing delivering ")
				assert.Equal(T, id1, id3, "Error in comparing delivering ")
			}
			wg.Done()
			break
		}
	}
}

// AssertCasualDelivering_one_to_many
//  Testing CASUAL ORDERED
func AssertCasualDelivering_one_to_many(T *testing.T, nm int, wg *sync.WaitGroup) {
	for {
		del1 := Node1_co.DelQueue
		del2 := Node2_co.DelQueue
		del3 := Node3_co.DelQueue
		if len(del1) == nm && len(del2) == nm && len(del3) == nm {
			log.Println("The test assertion..")
			for i := 0; i < nm; i++ {
				id1 := del1[i].MessageHeader["i"]
				id2 := del2[i].MessageHeader["i"]
				id3 := del3[i].MessageHeader["i"]
				assert.Equal(T, id1, id2, "Error in comparing delivering ")
				assert.Equal(T, id2, id3, "Error in comparing delivering ")
				assert.Equal(T, id1, id3, "Error in comparing delivering ")
			}
			wg.Done()
			break
		}

	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
var MA = basic.NewMessage(
	map[string]string{"c": "MA"},
	[]byte("MA"))
var MB = basic.NewMessage(
	map[string]string{"c": "MB"},
	[]byte("MB"))
var MC = basic.NewMessage(
	map[string]string{"c": "MC"},
	[]byte("MC"))
var MD = basic.NewMessage(
	map[string]string{"c": "MD"},
	[]byte("MD"))
var ME = basic.NewMessage(
	map[string]string{"c": "ME"},
	[]byte("ME"))
var MF = basic.NewMessage(
	map[string]string{"c": "MF"},
	[]byte("MF"))

var NODE1MSG = []basic.Message{MA, MB}
var NODE2MSG = []basic.Message{MC, MD}
var NODE3MSG = []basic.Message{ME, MF}

// AssertCasualDelivering_Many_to_Many
//      - Nodo1 sends Message MA / MB
//      - Nodo2 sends Message MC / MD
//      - Nodo3 sends Message ME / MF
//Casualty :
//      - (a)  MA causes MC & ME
//      - (b)  MB causes MD
// In this case, MF is a concurrent message
//
// Compatible sequences example with casual ordered  :
//      1 . MA -> MB -> ME -> MC -> MD -> MF
//      2.  MA -> MC -> MB -> ME -> MD -> MF
//      3.  MA -> ME -> MC -> MB -> MF -> MD
// Example of non compatible sequence :
//      1.  MA -> ME -> MD -> MB -> MC -> MF ( because MD before MB)
func AssertCasualDelivering_Many_to_Many(T *testing.T, nm int, wg *sync.WaitGroup) {
	for {
		del1 := Node1_co.DelQueue
		del2 := Node2_co.DelQueue
		del3 := Node3_co.DelQueue
		if len(del1) == nm && len(del2) == nm && len(del3) == nm {
			log.Println("The test assertion..")
			comp := func(q []basic.Message) bool {
				pa, err := GetPosition(q, "MA")
				if err != nil {
					return false
				}
				pb, err := GetPosition(q, "MB")
				if err != nil {
					return false
				}
				pc, err := GetPosition(q, "MC")
				if err != nil {
					return false
				}
				pd, err := GetPosition(q, "MD")
				if err != nil {
					return false
				}
				pe, err := GetPosition(q, "ME")
				if err != nil {
					return false
				}
				if pa < pc && pa < pe && pb < pd {
					return true
				}
				return false
			}
			if assert.True(T, comp(del1), "Error in comparing del1") {
				T.Log("Node1 respects the sequential order")
			}
			if assert.True(T, comp(del2), "Error in comparing del2") {
				T.Log("Node2 respects the sequential order")
			}
			if assert.True(T, comp(del3), "Error in comparing del3") {
				T.Log("Node3 respects the sequential order")
			}
			wg.Done()
			break
		}
	}
}

func GetPosition(m []basic.Message, c string) (int, error) {
	for i := 0; i < len(m); i++ {
		if m[i].MessageHeader["c"] == c {
			return i, nil
		}
	}
	return -1, errors.New("Not found message")
}

func printSequence(q []basic.Message) {
	for i := 0; i < 6; i++ {
		fmt.Print("\033[u\033[K  ", q[i].MessageHeader["c"])
	}
	fmt.Print("\n")
}

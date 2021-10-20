package main

import (
	"bufio"
	"github.com/msalvati1997/b1multicasting/pkg/basic"
	client "github.com/msalvati1997/b1multicasting/pkg/basic/client"
	"log"
	"os"
)

func main() {
	//Client connection to the server
	serverAddr := ":8080"
	Id := "myId"
	delay := 1 //delay for sending operations

	conn, err := client.Connect(serverAddr, delay)
	if err != nil {
		log.Println("Error in connecting to the server")
	}
	for {
		log.Println("Input : ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			log.Println("Input : ")
			text := scanner.Bytes()
			ch := make(chan bool, 1)
			err := conn.Send(Id, basic.NewMessage(make(map[string]string), text), &ch)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}

package main

import (
	"b1multicasting/pkg/basic"
	"bufio"
	client "b1multicasting/pkg/basic/client"
	"log"
	"os"
)

func main() {
	//Client connection to the server
	serverAddr := ":8080"
	Id := "myId"
	delay := 1 //delay for sending operations

	conn, err := client.Connect(serverAddr, uint(delay))
	if err != nil {
		log.Println("Error in connecting to the server")
	}
    for {
		log.Println("Input : ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			log.Println("Input : ")
			text := scanner.Bytes()
		err := conn.Send(Id, basic.NewMessage(text))
		if (err!=nil) {
			log.Println(err.Error())
		}
	} }
}

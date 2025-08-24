package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/taimats/chample/client"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("command should in the following way:")
		fmt.Println("[command] <your name>")
		os.Exit(1)
	}
	name := os.Args[1]

	urlStr := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("Hello, %s!!\n", name)
	fmt.Println("You're joining the hub")

	send := make(chan string)
	cl := client.NewClient(conn, send)
	go cl.Read()
	go cl.Send()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		send <- line
	}
}

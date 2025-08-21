package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/gorilla/websocket"
	"github.com/taimats/chample/client"
)

func main() {
	conn, res, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("failed to Dial:", err)
	}
	defer conn.Close()
	defer res.Body.Close()

	var buf bytes.Buffer
	io.Copy(&buf, res.Body)
	fmt.Println("response Body:", buf.String())

	done := make(chan struct{})
	cl := client.NewClient(conn, "test", done)

	fmt.Println("start chatting")
	errch := client.Chatting(cl)

	cl.Send <- client.NewMessage("client", "sending my message")

	err = <-errch
	if err != nil {
		conn.Close()
		log.Fatal("chatting error: ", err)
	}
}

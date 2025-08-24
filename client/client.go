package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan string
}

func NewClient(conn *websocket.Conn, send chan string) *Client {
	return &Client{
		conn: conn,
		send: send,
	}
}

func (cl *Client) Read() {
	defer cl.conn.Close()
	for {
		_, msg, err := cl.conn.ReadMessage()
		if err != nil {
			cl.conn.Close()
			return
		}
		fmt.Println(string(msg))
	}
}

func (cl *Client) Send() {
	defer func() {
		close(cl.send)
		cl.conn.Close()
	}()
	for {
		msg := <-cl.send
		err := cl.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			cl.conn.Close()
			log.Printf("WriteMessage error: %s\n", err)
			return
		}
	}
}

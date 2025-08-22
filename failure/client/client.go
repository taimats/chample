package client

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name string

	conn *websocket.Conn
}

func NewClient(name string, conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		Name: name,
	}
}

func (cl *Client) Read() {
	for {
		msg := new(Message)
		err := cl.conn.ReadJSON(msg)
		if err != nil {
			cl.conn.Close()
			log.Printf("failed to ReadJSON: (error: %s)", err)
			return
		}
		fmt.Printf("%s>> %s", msg.From, msg.Text)
	}
}

func (cl *Client) Write(msg *Message) error {
	err := cl.conn.WriteJSON(msg)
	if err != nil {
		return err
	}
	return nil
}

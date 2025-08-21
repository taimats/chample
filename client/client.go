package client

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client[T any] struct {
	name string

	conn *websocket.Conn
	send chan *Message
	done chan T
}

func NewClient[T any](conn *websocket.Conn, name string, done chan T) *Client[T] {
	return &Client[T]{
		conn: conn,
		name: name,
		send: make(chan *Message, 1),
		done: done,
	}
}

func (c *Client[T]) read() error {
	for {
		select {
		case <-c.done:
			return nil
		default:
			msg := new(Message)
			err := c.conn.ReadJSON(msg)
			if err != nil {
				return err
			}
			fmt.Printf("%s>> %s", msg.From, msg.Text)
		}
	}
}

func (c *Client[T]) write() error {
	for {
		select {
		case <-c.done:
			return nil
		default:
			err := c.conn.WriteJSON(<-c.send)
			if err != nil {
				return err
			}
		}
	}
}

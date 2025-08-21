package client

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client[T any] struct {
	name string

	conn *websocket.Conn
	Send chan *Message
	done chan T
}

func NewClient[T any](conn *websocket.Conn, name string, done chan T) *Client[T] {
	return &Client[T]{
		conn: conn,
		name: name,
		Send: make(chan *Message, 1),
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
				c.conn.Close()
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
			err := c.conn.WriteJSON(<-c.Send)
			if err != nil {
				c.conn.Close()
				return err
			}
		}
	}
}

func Chatting[T any](cl *Client[T]) chan error {
	errch := make(chan error)
	go func() {
		err := cl.read()
		if err != nil {
			errch <- fmt.Errorf("read error: %w", err)
		} else {
			errch <- nil
		}
	}()
	go func() {
		err := cl.write()
		if err != nil {
			errch <- fmt.Errorf("write error: %w", err)
		} else {
			errch <- nil
		}
	}()
	return errch
}

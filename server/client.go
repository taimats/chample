package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	recieve chan *Message
}

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:     h,
		conn:    conn,
		recieve: make(chan *Message),
	}
}

func (cl *Client) ReadConn() {
	defer cl.conn.Close()
	for {
		msgType, msg, err := cl.conn.ReadMessage()
		if err != nil {
			cl.hub.unregister <- cl
			log.Println(err)
			return
		}
		cl.hub.broadcast <- NewMessage(msgType, string(msg))
	}
}

func (cl *Client) WriteConn() {
	defer func() {
		close(cl.recieve)
		cl.conn.Close()
	}()
	for {
		msg := <-cl.recieve
		err := cl.conn.WriteMessage(msg.Type, []byte(msg.Content))
		if err != nil {
			cl.hub.unregister <- cl
			log.Println(err)
			return
		}
	}
}

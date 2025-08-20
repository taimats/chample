package example

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func NewClient() *Client {
	return &Client{}
}

// websocket conn --> message --> hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("failed to ReadMessage: (error: %v)", err)
			}
			return
		}
		c.send <- msg
	}
}

// websocket conn <-- message <-- hub
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer c.conn.Close()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			defer w.Close()
			w.Write(msg)
			w.Write([]byte{'\n'})
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

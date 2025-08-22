package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/taimats/chample/failure/client"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var GlobHub = NewHub()

type Client struct {
	id   string
	conn *websocket.Conn
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		id:   id,
		conn: conn,
	}
}

func ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Printf("failed to gen uuid: (error: %s)\n", err)
		conn.Close()
	}
	cl := NewClient(id.String(), conn)
	GlobHub.register <- cl
	log.Printf("accepting a new client: (id: %s)", cl.id)

	for {
		msg := new(client.Message)
		err := conn.ReadJSON(msg)
		if err != nil {
			conn.Close()
			log.Println("failed to ReadJSON:", err)
			return
		}
		GlobHub._broadcast(msg)
	}
}

type Hub struct {
	clients map[*Client]struct{}

	register  chan *Client
	broadcast chan *client.Message
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]struct{}),
		register:  make(chan *Client),
		broadcast: make(chan *client.Message),
	}
}

func (h *Hub) Register(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			cl := <-h.register
			h.clients[cl] = struct{}{}
		}
	}
}

func (h *Hub) _broadcast(msg *client.Message) {
	for cl := range h.clients {
		err := cl.conn.WriteJSON(msg)
		if err != nil {
			log.Println("failed to WriteJSON:", err)
			cl.conn.Close()
			delete(h.clients, cl)
			continue
		}
	}
}

func (h *Hub) Broadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg := <-h.broadcast
			for cl := range h.clients {
				err := cl.conn.WriteJSON(msg)
				if err != nil {
					log.Println("failed to WriteJSON:", err)
					cl.conn.Close()
					delete(h.clients, cl)
					continue
				}
			}
		}
	}
}

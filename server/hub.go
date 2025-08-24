package server

import "log"

type Hub struct {
	clients    map[*Client]struct{}
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	defer func() {
		close(h.broadcast)
		close(h.register)
		close(h.unregister)
	}()
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
			log.Printf("accepting a new client: %v", c)
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				log.Printf("leaving the hub: %v", c)
			}
		case msg := <-h.broadcast:
			for cl := range h.clients {
				select {
				case cl.recieve <- msg:
				default:
					delete(h.clients, cl)
				}
			}
		}
	}
}

package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024, //1KB
	WriteBufferSize: 1024,
}

func ServeWebsocket(h *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	cl := NewClient(h, conn)
	h.register <- cl
	// defer func() {
	// 	h.unregister <- cl
	// }()

	go cl.ReadConn()
	cl.WriteConn()
}

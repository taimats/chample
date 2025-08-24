package main

import (
	"log"
	"net/http"

	"github.com/taimats/chample/server"
)

func main() {
	addr := ":8080"
	h := server.NewHub()
	go h.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWebsocket(h, w, r)
	})
	log.Printf("server listening on port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

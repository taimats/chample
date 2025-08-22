package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/taimats/chample/failure/server"
)

var addr = ":8080"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go server.GlobHub.Register(ctx)
	go server.GlobHub.Broadcast(ctx)

	http.HandleFunc("/ws", http.HandlerFunc(server.ServeWebsocket))

	fmt.Printf("server listening on port %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

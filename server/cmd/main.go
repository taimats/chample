package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/taimats/chample/client"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalf("failed to upgrade: error %s", err)
		}
		defer conn.Close()

		for {
			msg := new(client.Message)
			err := conn.ReadJSON(msg)
			if err != nil {
				log.Fatalf("failed to readJson: error %s", err)
			}
			fmt.Printf("%s>> %s", msg.From, msg.Text)

			reply := client.NewMessage("server", "Thank you for you kind message!")
			err = conn.WriteJSON(reply)
			if err != nil {
				log.Fatalf("failed to write message: %s", err)
			}
		}
	})
	fmt.Println("server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

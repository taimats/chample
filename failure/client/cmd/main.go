package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/taimats/chample/failure/client"
)

var prompt = "%s>>"

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("your user name needed in the following way:")
		fmt.Println("command [user name]")
		os.Exit(1)
	}
	name := args[1]

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("failed to Dial:", err)
	}
	defer conn.Close()

	cl := client.NewClient(name, conn)

	fmt.Println("start chatting...")
	fmt.Printf("%s joining the hub!\n", name)

	sc := bufio.NewScanner(os.Stdin)
	go input(sc, cl)
	cl.Read()
}

func input(sc *bufio.Scanner, cl *client.Client) {
	for sc.Scan() {
		fmt.Printf(prompt, cl.Name)
		line := sc.Text()
		text := parseLine(line)
		if text == "" {
			continue
		}
		msg := client.NewMessage(cl.Name, text)
		err := cl.Write(msg)
		if err != nil {
			log.Printf("failed to WriteJSON: (error: %s)", err)
			return
		}
	}
}

func parseLine(line string) string {
	segs := strings.Split(line, ">>")
	if len(segs) != 2 {
		return ""
	}
	text := segs[1]
	return text
}

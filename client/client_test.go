package client_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/taimats/chample/client"
)

func TestClientRead(t *testing.T) {
	testmsg := client.NewMessage("test", "test message")

	path := "/ws"
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		h := newWebsocketHandler(conn)
		h.write(testmsg)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	u := &url.URL{}
	urlStruct, err := u.Parse(srv.URL)
	if err != nil {
		t.Fatal("url Parse error:", err)
	}
	urlStruct.Scheme = "ws"
	urlStr, err := url.JoinPath(urlStruct.String(), path)
	if err != nil {
		t.Fatal("JoinPath error:", err)
	}

	done := make(chan struct{})
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		t.Fatal("Dial error:", err)
	}
	defer conn.Close()
	cl := client.NewClient(conn, "test", done)

	//Act
	fnc := func() {
		client.ClientRead(cl)
	}
	got := dumpStdout(t, fnc)

	want := fmt.Sprintf("%s>> %s", "test", "test message")
	if got != want {
		t.Errorf("Not equal: (got=%s, want=%s)", got, want)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketHandler struct {
	conn *websocket.Conn
}

func newWebsocketHandler(conn *websocket.Conn) *websocketHandler {
	return &websocketHandler{conn: conn}
}

func (s *websocketHandler) read() {
	for {
		msg := new(client.Message)
		err := s.conn.ReadJSON(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s>> %s", msg.From, msg.Text)
	}
}

func (s *websocketHandler) write(msg *client.Message) {
	err := s.conn.WriteJSON(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func dumpStdout(t *testing.T, fnc func()) string {
	t.Helper()

	backup := os.Stdout
	defer func() {
		os.Stdout = backup
	}()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fnc()
	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("failed to readfrom: error: %s", err)
	}
	return buf.String()
}

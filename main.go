package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	tmp "github.com/lijiaqigreat/wsrpcgo/protobuf"
)

var (
	addr    = flag.String("addr", ":8080", "http service address")
	cmdPath string
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message
	// from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period.
	// Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on
	// connection.
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

/*
haha
*/
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeContent(w, r, "", time.Time{}, strings.NewReader("home.html"))
}

func main() {
	flag.Parse()
	var err error
	if err != nil {
		log.Fatal(err)
	}
	roomServer := NewRoomServer(nil)
	twirpHandler := tmp.NewRoomServiceServer(roomServer, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHome)
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)
	mux.Handle("/ws", roomServer.GetHandler())
	fmt.Printf("now serving %s\n", *addr)
	log.Fatal(http.ListenAndServeTLS(*addr, "server.cert", "server.key", mux))
}

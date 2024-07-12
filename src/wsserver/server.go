package wsserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	NONE     = "None"
	USERNAME = "Username"
)

type Server struct {
	hub          *Room
	clients      map[*Client]bool
	broadcast    chan []byte
	member_count int
	// later want reference to other server IP's so I can send to them as well (preparation of distributed network)
}

func NewServer() *Server {
	return &Server{
		hub:          NewRoom("Hub", nil),
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		member_count: 0,
	}
}

func (ws *Server) Run() {
	ws.hub.run()
	for {
		select {
		case msg := <-ws.broadcast:
			fmt.Println("Broadcasting " + string(msg) + " to all connected channels!")
		default:
			fmt.Println("Broadcast channel is full!")
			return
		}
	}
}

func (ws *Server) Serve(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		username: "guest",
		conn:     conn,
		room:     ws.hub,
		send:     make(chan []byte, 256),
		prompt:   "None",
	}
	ws.clients[client] = true
	client.room.register <- client
	go client.read()
	go client.write()
	// should then prompt for username
}

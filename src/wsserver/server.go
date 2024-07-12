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
	leaving      chan *Client
	member_count int
	// later want reference to other server IP's so I can send to them as well (preparation of distributed network)
}

func NewServer() *Server {
	s := &Server{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		leaving:      make(chan *Client),
		member_count: 0,
	}
	s.hub = NewRoom("hub", nil, s)
	return s
}

func (ws *Server) Run() {
	ws.hub.run()
	for {
		select {
		case msg := <-ws.broadcast:
			fmt.Println("Broadcasting " + string(msg) + " to all connected channels!")
		case client := <-ws.leaving:
			delete(ws.clients, client)
			close(client.send)
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

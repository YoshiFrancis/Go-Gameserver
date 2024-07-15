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
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, be cautious with this in production
	},
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
	TCPSend      chan []byte
	TCPRead      chan []byte
	member_count int
	// later want reference to other server IP's so I can send to them as well (preparation of distributed network)
}

func NewServer() *Server {
	s := &Server{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte, 1024),
		leaving:      make(chan *Client, 20),
		TCPSend:      make(chan []byte, 1024),
		member_count: 0,
	}
	s.hub = NewRoom("hub", nil, s)
	return s
}

func (ws *Server) Run() {
	go ws.hub.run()
	for {
		select {
		case msg := <-ws.broadcast:
			fmt.Println("Broadcasting " + string(msg) + " to all connected channels!")
			for client := range ws.clients {
				client.send <- msg
			}
		case client := <-ws.leaving:
			delete(ws.clients, client)
			close(client.send)
		case message := <-ws.TCPRead:
			fmt.Println("MEssage received via tcpserver in websocket server")
			ws.broadcast <- message
		}
	}
}

func (ws *Server) Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User has connected!")
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
		prompt:   USERNAME,
	}
	ws.clients[client] = true
	client.room.register <- client
	go client.read()
	go client.write()
}

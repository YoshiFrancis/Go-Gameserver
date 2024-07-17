package wsserver

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/src/messages"
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

type WSServer struct {
}

type Server struct {
	hub          *Room
	clients      map[string]*Client
	rooms        map[int]*Room
	broadcast    chan []byte
	leaving      chan string
	TCPSend      chan []byte
	TCPRead      chan []byte
	member_count int
	roomIdCount  int
	// later want reference to other server IP's so I can send to them as well (preparation of distributed network)
}

func NewServer() *Server {
	s := &Server{
		clients:      make(map[string]*Client),
		rooms:        make(map[int]*Room),
		broadcast:    make(chan []byte, 1024),
		leaving:      make(chan string, 20),
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
			for _, client := range ws.clients {
				client.send <- msg
			}
		case client := <-ws.leaving:
			ws.clients[client].room.unregister <- client
			close(ws.clients[client].send)
			delete(ws.clients, client)
		case message := <-ws.TCPRead:
			flag, args := messages.Decode(message)
			if flag == '-' {
				go ws.handleCommand(args)
			} else if flag == '+' {
				roomId, _ := strconv.Atoi(args[1])
				go ws.rooms[roomId].handleCommand(args)
			} else if flag == '*' {
				username := args[1]
				go ws.clients[username].handleCommand(args)
			}
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

	go client.read()
	go client.write()
}

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

type WSServer struct {
	clients   map[string]*Client
	broadcast chan []byte
	leaving   chan string
	register  chan string
	requests  chan []byte
	serverId  int
}

func NewWSServer() *WSServer {
	return &WSServer{
		clients:   make(map[string]*Client),
		broadcast: make(chan []byte, 1024),
		leaving:   make(chan string, 20),
		register:  make(chan string, 12),
		requests:  make(chan []byte, 1024),
	}
}

func (ws *WSServer) Run() {
	defer ws.shutdown()
	for {
		select {
		case msg := <-ws.broadcast:
			for _, client := range ws.clients {
				client.send <- msg
			}
		case client := <-ws.leaving:
			close(ws.clients[client].send)
			delete(ws.clients, client)
			// signal Leader if have one
		}
	}
}

func (ws *WSServer) shutdown() {
	close(ws.broadcast)
	close(ws.leaving)
}

func (ws *WSServer) Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Println("User has connected!")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go ws.getUsername(conn)
}

func (ws *WSServer) getUsername(conn *websocket.Conn) {
	client := &Client{
		username: "",
		conn:     conn,
		send:     make(chan []byte, 256),
		server:   ws,
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Client is going to stop reading!")
		conn.Close()
		return
	}

	client.username = string(message)
	go client.read()
	go client.write()

}

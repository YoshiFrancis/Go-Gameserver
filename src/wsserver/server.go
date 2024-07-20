package wsserver

import (
	"fmt"
	"log"
	"net/http"

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
	Clients   map[string]*Client
	broadcast chan []byte
	leaving   chan string
	register  chan string
	requests  chan []byte
	serverId  int
}

func NewWSServer(requests chan []byte) *WSServer {
	return &WSServer{
		Clients:   make(map[string]*Client),
		broadcast: make(chan []byte, 1024),
		leaving:   make(chan string, 20),
		register:  make(chan string, 12),
		requests:  requests,
		serverId:  -1,
	}
}

func (ws *WSServer) Run() {
	for {
		select {
		case msg := <-ws.broadcast:
			for _, client := range ws.Clients {
				client.Send <- msg
			}
		case client := <-ws.leaving:
			close(ws.Clients[client].Send)
			delete(ws.Clients, client)
		}
	}
}

func (ws *WSServer) Shutdown() {
	fmt.Println("wsserver shutting down")
	close(ws.broadcast)
	close(ws.leaving)
	close(ws.register)
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
		Send:     make(chan []byte, 256),
		server:   ws,
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Client is going to stop reading!")
		conn.Close()
		return
	}

	client.username = string(message)
	fmt.Println("New client!", client.username)
	register_msg := messages.ServerJoinUser(client.username, ws.serverId)
	ws.requests <- []byte(register_msg)
	go client.read()
	go client.write()

}

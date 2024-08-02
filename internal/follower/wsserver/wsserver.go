package wsserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

var usernameTemplate *template.Template
var hubTemplate *template.Template

func init() {
	hubTemplate = template.Must(template.ParseFiles("../internal/follower/wsserver/hub.html"))
	usernameTemplate = template.Must(template.ParseFiles("../internal/follower/wsserver/username.html"))
}

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

type Username struct {
	Username string `json:"username"`
}

type WSServer struct {
	Clients    map[string]*Client
	broadcast  chan []byte
	unregister chan string
	register   chan string
	TCPfrom    chan []byte
	TCPto      chan []byte
	ServerId   string
}

func NewWSServer() *WSServer {
	return &WSServer{
		Clients:    make(map[string]*Client),
		broadcast:  make(chan []byte),
		unregister: make(chan string),
		register:   make(chan string),
	}
}

func (ws *WSServer) Username(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var username Username
		err := conn.ReadJSON(&username)
		if err != nil {
			log.Println("Read error of username:", err)
			break
		}

		log.Printf("Received username: %s", username.Username)
	}
}

func (ws *WSServer) Run() {
	for {
		select {
		case msg := <-ws.broadcast:
			for _, client := range ws.Clients {
				client.Send <- msg
			}
		case client := <-ws.unregister:
			close(ws.Clients[client].Send)
			delete(ws.Clients, client)
		}
	}
}

func (ws *WSServer) Shutdown() {
	fmt.Println("wsserver shutting down")
	close(ws.broadcast)
	close(ws.unregister)
	close(ws.register)
}

func (ws *WSServer) Home(w http.ResponseWriter, r *http.Request) {
	err := usernameTemplate.Execute(w, struct {
		WebsocketHost string
	}{
		WebsocketHost: "ws://" + r.Host + "/username",
	})

	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		fmt.Println("Error rendering template:", err)
		return
	}
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

	var username Username
	json.Unmarshal(message, &username)
	client.username = username.Username
	fmt.Println("New client!", client.username)
	register_msg := messages.ServerJoinUser(client.username, -1)
	ws.TCPto <- []byte(register_msg)

	go client.read()
	go client.write()
}

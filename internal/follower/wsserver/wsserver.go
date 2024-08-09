package wsserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/internal/containers"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

var indexTemplate *template.Template
var usernameTemplate *template.Template
var hubTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.ParseFiles("../web/index.html"))
	usernameTemplate = template.Must(template.ParseFiles("../internal/templates/username.html"))
	hubTemplate = template.Must(template.ParseFiles("../internal/templates/hub.html"))
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

type WSServer struct {
	Clients    *containers.Storage[string, *Client]
	broadcast  chan []byte
	unregister chan *Client
	register   chan *Client
	TCPfrom    chan []byte
	TCPto      chan []byte
	done       chan bool
	ServerId   string
}

func NewWSServer(done chan bool) *WSServer {
	return &WSServer{
		Clients:    containers.NewStorage[string, *Client](),
		broadcast:  make(chan []byte, 24),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		done:       done,
	}
}

type Username struct {
	Username string `json:"username"`
}

func (ws *WSServer) Home(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}

	promptUsernameTemp := containers.RenderTemplate(usernameTemplate, struct{}{})
	conn.WriteMessage(websocket.TextMessage, promptUsernameTemp)

	// Read message from browser
	_, p, err := conn.ReadMessage()

	if err != nil {
		log.Println("error reading: ", err)
		return
	}

	var username Username
	err = json.Unmarshal(p, &username)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("New user: ", username.Username)

	c := NewClient(username.Username, conn, ws)
	register_msg := messages.ServerRegisterUser(username.Username, -1)
	ws.TCPto <- []byte(register_msg)
	hubTemp := containers.RenderTemplate(hubTemplate, struct{ Username string }{Username: username.Username})
	c.send <- hubTemp
}

func (ws *WSServer) Ping(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
		return
	}
	_, p, err := conn.ReadMessage()

	if err != nil {
		log.Println("error reading: ", err)
		return
	}
	flag, ping := messages.Decode(p)
	if flag == '!' && ping[0] == "PING" {
		err := conn.WriteMessage(websocket.TextMessage, []byte(messages.Pong()))
		if err != nil {
			fmt.Println("Error writing pong message: ", err)
		}
	}
}

func (ws *WSServer) Run() {
	for {
		select {
		case msg := <-ws.broadcast:
			for _, client := range ws.Clients.Values() {
				client.send <- msg
			}
		case client := <-ws.register:
			fmt.Println("Registering client!", client.username)
			ws.Clients.Set(client.username, client)
		case client := <-ws.unregister:
			if _, ok := ws.Clients.Get(client.username); ok {
				close(client.send)
				ws.Clients.Delete(client.username)
			}
		case from := <-ws.TCPfrom:
			ws.broadcast <- from
		}
	}
}

func (ws *WSServer) Shutdown() {
	fmt.Println("wsserver shutting down")
	close(ws.broadcast)
	close(ws.unregister)
	close(ws.register)
	if ws.done != nil {
		ws.done <- true
	}
}

func (ws *WSServer) Index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		fmt.Println("Error rendering username template:", err)
		return
	}
}

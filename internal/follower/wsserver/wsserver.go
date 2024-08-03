package wsserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/internal/leader/storage"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

var indexTemplate *template.Template
var usernameTemplate *template.Template

func init() {
	indexTemplate = template.Must(template.ParseFiles("../web/index.html"))
	usernameTemplate = template.Must(template.ParseFiles("../internal/templates/username.html"))
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
	Clients    *storage.Storage[string, *Client]
	broadcast  chan []byte
	unregister chan *Client
	register   chan *Client
	TCPfrom    chan []byte
	TCPto      chan []byte
	ServerId   string
}

func NewWSServer() *WSServer {
	return &WSServer{
		Clients:    storage.NewStorage[string, *Client](),
		broadcast:  make(chan []byte),
		unregister: make(chan *Client),
		register:   make(chan *Client),
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

	promptUsernameTemp := renderTemplate(struct{}{})
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

	NewClient(username.Username, conn, ws)
	register_msg := messages.ServerJoinUser(username.Username, -1)
	ws.TCPto <- []byte(register_msg)

}

func (ws *WSServer) Run() {
	for {
		select {
		case msg := <-ws.broadcast:
			for _, client := range ws.Clients.Values() {
				client.Send <- msg
			}
		case client := <-ws.register:
			ws.Clients.Set(client.username, client)
		case client := <-ws.unregister:
			close(client.Send)
			ws.Clients.Delete(client.username)
		}
	}
}

func (ws *WSServer) Shutdown() {
	fmt.Println("wsserver shutting down")
	close(ws.broadcast)
	close(ws.unregister)
	close(ws.register)
}

func (ws *WSServer) Index(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		fmt.Println("Error rendering username template:", err)
		return
	}
}

func renderTemplate[T any](dataStruct T) []byte {
	var templateBuffer bytes.Buffer
	usernameTemplate.Execute(&templateBuffer, dataStruct)
	return templateBuffer.Bytes()
}

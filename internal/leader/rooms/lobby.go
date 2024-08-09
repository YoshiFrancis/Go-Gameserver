package rooms

import (
	"html/template"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
)

var lobbyTemplate *template.Template

func init() {
	lobbyTemplate = template.Must(template.ParseFiles("../internal/templates/lobby.html"))
}

type Lobby struct {
	users    *containers.Storage[string, User]
	roomId   int
	registry chan User
	prevRoom Room
	title    string
	msgHist  *containers.Queue[Message]
	creator  string
}

func NewLobby(id int, prevRoom Room, title, creator string) *Lobby {
	return &Lobby{
		users:    containers.NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
		prevRoom: prevRoom,
		title:    title,
		msgHist:  containers.NewQueue[Message](20),
		creator:  creator,
	}
}

func (l *Lobby) Join(user User) []byte {
	user.room.Leave(user)
	l.users.Set(user.username, user)
	return containers.RenderTemplate(lobbyTemplate, struct {
		LobbyTitle      string
		Username        string
		CreatorUsername string
	}{
		LobbyTitle:      l.title,
		Username:        user.username,
		CreatorUsername: l.creator,
	})
}

func (l *Lobby) Leave(user User) {
	l.users.Delete(user.username)
	l.prevRoom.Join(user)
}

func (l *Lobby) Broadcast(sender, message string) string {
	l.msgHist.Enqueue(Message{sender, message})
	return l.getHTMXMessages()
}

func (l *Lobby) HandleMessage(message string, sender string) {
	// --------------------------- lobby handle message ---------------------------
}

func (l *Lobby) GetInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the lobby."
}

func (l *Lobby) GetName() string {
	return l.title
}

func (l *Lobby) getHTMXMessages() string {
	messages := l.msgHist.Items()
	htmx := "<div id=\"chat-room\" hx-swap=\"outerHTML\"><ul>"
	for _, message := range messages {
		htmx += "<li>" + message.username + ": " + message.text + "</li>"
	}
	htmx += "</ul></div>"
	return htmx
}

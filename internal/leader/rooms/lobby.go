package rooms

import (
	"fmt"
	"html/template"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
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
	app      string
}

func NewLobby(id int, prevRoom Room, title, creator string) *Lobby {
	fmt.Println("Creating new lobby: ", title)
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

func (l *Lobby) Join(user *User) (leavingTmpl, joiningTmpl []byte) {

	leavingTmpl = user.room.Leave(user)

	l.users.Set(user.username, *user)
	user.room = l
	fmt.Println(user.username + " has joined " + user.room.GetName())
	l.msgHist.Enqueue(Message{
		Username: "Server",
		Text:     user.username + " has joined!",
	})

	joiningTmpl = containers.RenderTemplate(lobbyTemplate, struct {
		LobbyTitle      string
		Username        string
		CreatorUsername string
		Participants    []string
		Messages        []Message
	}{
		LobbyTitle:      l.title,
		Username:        user.username,
		CreatorUsername: l.creator,
		Participants:    l.users.Keys(),
		Messages:        l.msgHist.Items(),
	})

	return leavingTmpl, joiningTmpl
}

func (l *Lobby) Leave(user *User) []byte {
	fmt.Println(user.username + " is leaving " + l.title)
	l.users.Delete(user.username)
	user.room = nil

	leavingTmpl := containers.RenderTemplate(lobbyTemplate, struct {
		LobbyTitle      string
		CreatorUsername string
		Participants    []string
		Messages        []Message
	}{
		LobbyTitle:      l.title,
		CreatorUsername: l.creator,
		Participants:    l.users.Keys(),
		Messages:        l.msgHist.Items(),
	})

	return leavingTmpl
}

func (l *Lobby) BroadcastMessage(sender, message string) string {
	l.msgHist.Enqueue(Message{sender, message})
	broadcastMsg := messages.LeaderRoomBroadcast(l.getHTMXMessages(), l.users.Keys())
	return broadcastMsg
}

func (l *Lobby) BroadcastTemplate(tmpl string) string {
	return messages.LeaderRoomBroadcast(tmpl, l.users.Keys())
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

func (l *Lobby) LeavingMessage(leavingUser string) string {
	return messages.LeaderRoomBroadcast(leavingUser+" has left!", l.users.Keys())
}

func (l *Lobby) JoiningMessage(joiningUser string) string {
	return messages.LeaderRoomBroadcast(joiningUser+" has left!", l.users.Keys())
}

func (l *Lobby) getHTMXMessages() string {
	messages := l.msgHist.Items()
	htmx := "<div id=\"chat-room\" hx-swap=\"outerHTML\"><ul>"
	for _, message := range messages {
		htmx += "<li>" + message.Username + ": " + message.Text + "</li>"
	}
	htmx += "</ul></div>"
	return htmx
}

func (l *Lobby) GetApp() string {
	return l.app
}

func (l *Lobby) GetUsers() []string {
	return l.users.Keys()
}

func (l *Lobby) SetApp(appName string) {
	l.app = appName
}

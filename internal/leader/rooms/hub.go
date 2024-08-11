package rooms

import (
	"fmt"
	"html/template"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

var hubTemplate *template.Template

func init() {
	hubTemplate = template.Must(template.ParseFiles("../internal/templates/hub.html"))
}

type Hub struct {
	users    *containers.Storage[string, User]
	roomId   int
	registry chan User
	msgHist  *containers.Queue[Message]
	app      string
}

func NewHub(id int) *Hub {
	return &Hub{
		users:    containers.NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
		msgHist:  containers.NewQueue[Message](20),
		app:      "",
	}
}

func (h *Hub) Join(user *User) (leavingTmpl, joiningTmpl []byte) {
	leavingTmpl = user.room.Leave(user)
	fmt.Println(user.username + " is joining the hub!")
	h.users.Set(user.username, *user)
	user.room = h
	fmt.Println("People in hub: ", h.users.Keys())
	fmt.Println(h.msgHist.Items())
	joiningTmpl = containers.RenderTemplate(hubTemplate, struct {
		Participants []string
		Messages     []Message
	}{
		Participants: h.users.Keys(),
		Messages:     h.msgHist.Items(),
	})

	return leavingTmpl, joiningTmpl
}

func (h *Hub) Leave(user *User) []byte {
	fmt.Println(user.username + " is leaving the hub!")
	h.users.Delete(user.username)
	fmt.Println("People in hub: ", h.users.Keys())
	user.room = nil

	leavingTmpl := containers.RenderTemplate(hubTemplate, struct {
		Participants []string
		Messages     []Message
	}{
		Participants: h.users.Keys(),
		Messages:     h.msgHist.Items(),
	})

	return leavingTmpl

}

func (h *Hub) BroadcastMessage(sender, message string) string {
	h.msgHist.Enqueue(Message{sender, message})
	broadcastMsg := messages.LeaderRoomBroadcast(h.getHTMXMessages(), h.users.Keys())
	return broadcastMsg
}

func (h *Hub) BroadcastTemplate(tmpl string) string {
	return messages.LeaderRoomBroadcast(tmpl, h.users.Keys())
}

func (h *Hub) GetInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the hub. This is the default area where all users are sent to on joining."
}

func (h *Hub) GetId() int {
	return h.roomId
}

func (h *Hub) LeavingMessage(leavingUser string) string {
	return messages.LeaderRoomBroadcast(leavingUser+" has left!", h.users.Keys())
}

func (h *Hub) JoiningMessage(joiningUser string) string {
	return messages.LeaderRoomBroadcast(joiningUser+" has left!", h.users.Keys())
}

func (h *Hub) getHTMXMessages() string {
	messages := h.msgHist.Items()
	htmx := "<div id=\"chat-room\" hx-swap=\"outerHTML\"><ul>"
	for _, message := range messages {
		htmx += "<li>" + message.Username + ": " + message.Text + "</li>"
	}
	htmx += "</ul></div>"
	return htmx
}

func (h *Hub) GetName() string {
	return "Hub"
}

func (h *Hub) GetApp() string {
	return h.app // should return ""
}

func (h *Hub) GetUsers() []string {
	return h.users.Keys()
}

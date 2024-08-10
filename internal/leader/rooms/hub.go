package rooms

import (
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
}

func NewHub(id int) *Hub {
	return &Hub{
		users:    containers.NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
		msgHist:  containers.NewQueue[Message](20),
	}
}

func (h *Hub) Join(user User) []byte {
	user.room.Leave(user)
	h.users.Set(user.username, user)
	return containers.RenderTemplate(hubTemplate, struct{ Username string }{Username: user.username})
}

func (h *Hub) Leave(user User) {
	h.users.Delete(user.username)
}

func (h *Hub) Broadcast(sender, message string) string {
	h.msgHist.Enqueue(Message{sender, message})
	broadcastMsg := messages.LeaderRoomBroadcast(h.getHTMXMessages(), h.users.Keys())
	return broadcastMsg
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
		htmx += "<li>" + message.username + ": " + message.text + "</li>"
	}
	htmx += "</ul></div>"
	return htmx
}

func (h *Hub) GetName() string {
	return "Hub"
}

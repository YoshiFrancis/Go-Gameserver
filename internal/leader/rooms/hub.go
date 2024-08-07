package rooms

import (
	"github.com/yoshifrancis/go-gameserver/internal/containers"
)

type Message struct {
	username string
	text     string
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
		msgHist:  &containers.Queue[Message]{},
	}
}

func (h *Hub) Join(user User) {
	user.room.Leave(user)
	h.users.Set(user.username, user)
}

func (h *Hub) Leave(user User) {
	h.users.Delete(user.username)
}

func (h *Hub) Broadcast(sender, message string) string {
	h.msgHist.Enqueue(Message{sender, message})
	return h.getHTMXMessages()
}

func (h *Hub) GetInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the hub. This is the default area where all users are sent to on joining."
}

func (h *Hub) GetId() int {
	return h.roomId
}

func (h *Hub) getUserStorage() *containers.Storage[string, User] {
	return h.users
}

func (h *Hub) getHTMXMessages() string {
	messages := h.msgHist.Items()
	htmx := "<ul>"
	for _, message := range messages {
		htmx += "<li>" + message.username + ": " + message.text + "<\\li>"
	}
	return htmx
}

package rooms

import (
	"fmt"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
)

type Hub struct {
	users    *containers.Storage[string, User]
	roomId   int
	registry chan User
}

func NewHub(id int) *Hub {
	return &Hub{
		users:    containers.NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
	}
}

func (h *Hub) Join(user User) {
	user.room.Leave(user)
	h.users.Set(user.username, user)
}

func (h *Hub) Leave(user User) {
	h.users.Delete(user.username)
}

func (h *Hub) DeliverAll(message string) {
	for user := range h.users.Values() {
		// --------------------------- send user message ---------------------------
		fmt.Println("Message for ", user)
	}
}

func (h *Hub) HandleMessage(message string, sender string) {
	// --------------------------- hub handle message ---------------------------
}

func (h *Hub) GetInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the hub. This is the default area where all users are sent to on joining."
}

func (h *Hub) getUserStorage() *containers.Storage[string, User] {
	return h.users
}

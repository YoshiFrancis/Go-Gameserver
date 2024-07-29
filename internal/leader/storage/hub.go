package storage

import "fmt"

type Hub struct {
	users    *Storage[string, User]
	roomId   int
	registry chan User
}

func NewHub(id int) *Hub {
	return &Hub{
		users:    NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
	}
}

func (h *Hub) join(user User) {
	user.room.leave(user)
	h.users.Set(user.username, user)
}

func (h *Hub) leave(user User) {
	h.users.Delete(user.username)
}

func (h *Hub) deliverAll(message string) {
	for user := range h.users.Values() {
		// --------------------------- send user message ---------------------------
		fmt.Println("Message for ", user)
	}
}

func (h *Hub) handleMessage(message string, sender string) {
	// --------------------------- hub handle message ---------------------------
}

func (h *Hub) getInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the hub. This is the default area where all users are sent to on joining."
}

func (h *Hub) getUserStorage() *Storage[string, User] {
	return h.users
}

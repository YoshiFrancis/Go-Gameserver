package storage

import "fmt"

type Lobby struct {
	users    *Storage[string, User]
	roomId   int
	registry chan User
	prevRoom Room
}

func NewLobby(id int, prevRoom Room) *Lobby {
	return &Lobby{
		users:    NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
		prevRoom: prevRoom,
	}
}

func (l *Lobby) join(user User) {
	user.room.leave(user)
	l.users.Set(user.username, user)
}

func (l *Lobby) leave(user User) {
	l.users.Delete(user.username)
	l.prevRoom.join(user)
}

func (l *Lobby) deliverAll(message string) {
	for user := range l.users.Values() {
		// --------------------------- send user message ---------------------------
		fmt.Println("Message for ", user)
	}
}

func (l *Lobby) handleMessage(message string, sender string) {
	// --------------------------- lobby handle message ---------------------------
}

func (l *Lobby) getInfo() string {
	// --------------------------- lobby info message ---------------------------
	return "This is the lobby."
}

func (l *Lobby) getUserStorage() *Storage[string, User] {
	return l.users
}

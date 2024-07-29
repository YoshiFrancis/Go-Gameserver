package storage

type Lobby struct {
	users    *Storage[string, User]
	roomId   int
	registry chan User
}

func NewLobby(id int) *Lobby {
	return &Lobby{
		users:    NewStorage[string, User](),
		roomId:   id,
		registry: make(chan User, 8),
	}
}

func (h *Lobby) join(username string) {

}

func (h *Lobby) leave(username string) {

}

func (h *Lobby) deliverAll(message string) {

}

func (h *Lobby) handleMessage(message string, sender string) {

}

func (h *Lobby) getInfo() string {
	return "This is the lobby."
}

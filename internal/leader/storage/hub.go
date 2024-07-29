package storage

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

func (h *Hub) join(username string) {

}

func (h *Hub) leave(username string) {

}

func (h *Hub) deliverAll(message string) {

}

func (h *Hub) handleMessage(message string, sender string) {

}

func (h *Hub) getInfo() string {
	return "This is the hub. This is the default area where all users are sent to on joining."
}

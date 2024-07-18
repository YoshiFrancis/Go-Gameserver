package leaderserver

import "fmt"

type Hub struct {
	hubId        int
	users        map[string]*User
	member_count int
	broadcast    chan []byte
	register     chan *User
	unregister   chan *User
}

func NewHub(id int) *Hub {
	return &Hub{
		hubId:        id,
		users:        make(map[string]*User),
		member_count: 0,
		broadcast:    make(chan []byte, 156),
		register:     make(chan *User, 156),
		unregister:   make(chan *User, 156),
	}
}

func (h *Hub) close() {
	close(h.broadcast)
}

func (h *Hub) run() {
	defer h.close()

	for {
		select {
		case user := <-h.register:
			h.users[user.username] = user
			user.roomId = h.hubId
			h.member_count++
			fmt.Println(user.username + " registered into hub")
		case user := <-h.unregister:
			delete(h.users, user.username)
		case msg := <-h.broadcast:
			for _, user := range h.users {
				user.send(msg)
			}
		}
	}
}

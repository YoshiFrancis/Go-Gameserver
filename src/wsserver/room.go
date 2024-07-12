package wsserver

import "fmt"

type Room struct {
	clients      map[*Client]bool
	parentRoom   *Room
	title        string
	member_count int
	register     chan *Client
	unregister   chan *Client
	messages     chan string
	server       *Server
}

func NewRoom(title string, parentRoom *Room, server *Server) *Room {
	if title == "" {
		title = "Lobby"
	}

	return &Room{
		clients:      make(map[*Client]bool),
		parentRoom:   parentRoom,
		title:        title,
		member_count: 0,
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		messages:     make(chan string),
		server:       server,
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
		case client := <-r.unregister:
			delete(r.clients, client)
			r.member_count--
			if r.parentRoom != nil {
				r.parentRoom.register <- client
			} else {
				r.server.leaving <- client
			}
		case message := <-r.messages:
			fmt.Println(message)
		}
	}

}

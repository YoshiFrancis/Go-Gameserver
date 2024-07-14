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
		register:     make(chan *Client, 28),
		unregister:   make(chan *Client, 28),
		messages:     make(chan string, 1024),
		server:       server,
	}
}

func (r *Room) run() {
	fmt.Println("Room is running!")
	for {
		select {
		case client := <-r.register:
			fmt.Println("user has registered!")
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
			fmt.Println("Message received: ", message)
		}
	}

}

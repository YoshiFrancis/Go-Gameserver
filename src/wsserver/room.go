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
	shutdown     chan bool
	server       *Server
	roomId       int
}

func NewRoom(title string, parentRoom *Room, server *Server) *Room {
	if title == "" {
		title = "Lobby"
	}
	server.roomIdCount++
	return &Room{
		clients:      make(map[*Client]bool),
		parentRoom:   parentRoom,
		title:        title,
		member_count: 0,
		register:     make(chan *Client, 10),
		unregister:   make(chan *Client, 10),
		messages:     make(chan string, 1024),
		shutdown:     make(chan bool, 1),
		server:       server,
		roomId:       server.roomIdCount,
	}
}

func (r *Room) run() {
	fmt.Println("Room is running!")
	defer r.close()
	for {
		select {
		case client := <-r.register:
			fmt.Println("user has registered!")
			r.clients[client] = true
			r.member_count++
			client.room = r
		case client := <-r.unregister:
			delete(r.clients, client)
			r.member_count--
			if r.parentRoom != nil {
				r.parentRoom.register <- client
			} else {
				r.server.leaving <- client
			}
			if r.member_count == -1 && r.title != "hub" {
				r.shutdown <- true
			}
		case message := <-r.messages:
			fmt.Println("Message received: ", message)
			r.server.TCPSend <- []byte(message)
		case <-r.shutdown:
			return
		}
	}
}

func getRoom(s *Server, roomId int) (*Room, bool) {
	new_room, ok := s.rooms[roomId]
	return new_room, ok
}

func (r *Room) close() {
	close(r.messages)
	close(r.register)
	close(r.unregister)
	close(r.shutdown)
	for client := range r.clients {
		client.switchRoom(r.parentRoom)
	}
	delete(r.server.rooms, r.roomId)
}

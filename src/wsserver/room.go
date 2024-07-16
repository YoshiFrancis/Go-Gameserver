package wsserver

import (
	"fmt"
)

type Room struct {
	clients      map[string]*Client
	parentRoom   *Room
	title        string
	member_count int
	register     chan *Client
	unregister   chan string
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
	new_room := &Room{
		clients:      make(map[string]*Client),
		parentRoom:   parentRoom,
		title:        title,
		member_count: 0,
		register:     make(chan *Client, 10),
		unregister:   make(chan string, 10),
		messages:     make(chan string, 1024),
		shutdown:     make(chan bool, 1),
		server:       server,
		roomId:       server.roomIdCount,
	}
	server.rooms[server.roomIdCount] = new_room
	return new_room
}

func (r *Room) run() {
	fmt.Println("Room is running!")
	defer r.close()
	for {
		select {
		case client := <-r.register:
			r.clients[client.username] = client
			r.member_count++
			client.room = r
			fmt.Println("user has registered to: ", r.title)
		case client := <-r.unregister:
			delete(r.clients, client)
			r.member_count--
			if r.member_count == -1 && r.title != "hub" {
				r.shutdown <- true
			}

		case message := <-r.messages:
			fmt.Println("Message received in " + r.title + ": " + message)
			// just broadcast for now
			for _, c := range r.clients {
				c.send <- []byte(message)
			}
			r.server.TCPSend <- []byte(message)
		case <-r.shutdown:
			return
		}
	}
}

func getRoom(s *Server, roomId int) (*Room, bool) {
	fmt.Println(s.rooms)
	for id, r := range s.rooms {
		fmt.Println(r.title+" with an id of", id)
	}
	new_room, ok := s.rooms[roomId]
	return new_room, ok
}

func (r *Room) close() {
	close(r.messages)
	close(r.register)
	close(r.unregister)
	close(r.shutdown)
	for _, client := range r.clients {
		client.switchRoom(r.parentRoom)
	}
	delete(r.server.rooms, r.roomId)
}

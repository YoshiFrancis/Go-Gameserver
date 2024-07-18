package leaderserver

import "fmt"

type Lobby struct {
	lobbyId      int
	users        map[string]*User
	member_count int
	broadcast    chan []byte
	register     chan *User
	unregister   chan *User
}

func NewLobby(id int) *Lobby {
	return &Lobby{
		lobbyId:      id,
		users:        make(map[string]*User),
		member_count: 0,
		broadcast:    make(chan []byte, 156),
		register:     make(chan *User, 4),
		unregister:   make(chan *User, 4),
	}
}

func (lobby *Lobby) close() {
	close(lobby.broadcast)
}

func (lobby *Lobby) run() {
	defer lobby.close()

	for {
		select {
		case user := <-lobby.register:
			lobby.users[user.username] = user
			user.roomId = lobby.lobbyId
			fmt.Println("User has joined lobby with id", lobby.lobbyId)
		case user := <-lobby.unregister:
			delete(lobby.users, user.username)
		case msg := <-lobby.broadcast:
			for _, user := range lobby.users {
				user.send(msg)
			}
		}
	}
}

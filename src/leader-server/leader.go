package leaderserver

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/yoshifrancis/go-gameserver/src/messages"
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

type Leader struct {
	WSServer    *wsserver.WSServer   // access to websocket server
	TCPServer   *tcpserver.TCPServer // access to tcp server
	WSrequests  chan []byte          // channel of incoming requests from websocket server
	TCPrequests chan []byte          // channel of incoming requests from tcp server
	isLeader    bool                 // indicate if leader or not
	idGen       func() int           // used to generate ids, get function from a closure
	hub         *Hub                 // room title, room id
	lobbies     map[int]*Lobby       // lobby id, Lobby pointer
	Users       map[string]*User     // username, user pointer
	mutex       sync.Mutex
}

func NewLeader() *Leader {
	return &Leader{
		Users:       make(map[string]*User),
		hub:         NewHub(0),
		lobbies:     make(map[int]*Lobby),
		WSrequests:  make(chan []byte, 1024),
		TCPrequests: make(chan []byte, 1024),
		idGen:       idGenerator(0),
		isLeader:    true,
		mutex:       sync.Mutex{},
	}
}

func (l *Leader) Run() {
	defer l.shutdown()

	for {
		select {
		case req := <-l.WSrequests:
			if l.isLeader {
				flag, args := messages.Decode(req)
				message := l.handleArgs(flag, args)
				fmt.Println("Received message from WSSS:", string(message))
				l.TCPServer.Broadcast <- []byte(message)
			} else {
				fmt.Println("Received from WSServer", string(req))
				l.TCPServer.Broadcast <- req // should give to leader
			}
		case req := <-l.TCPrequests:
			fmt.Println("Received message from tcp server")
			flag, args := messages.Decode(req)
			if l.isLeader {
				message := l.handleArgs(flag, args)
				fmt.Println(message)
				l.TCPServer.Broadcast <- []byte(message)
			} else {
				if args[0] == "BROADCAST" {
					if flag == '-' {
						if hubId, _ := strconv.Atoi(args[1]); hubId != l.hub.hubId {
							fmt.Println("Given invalid hub id!")
							continue
						}
						l.hub.broadcast <- []byte(args[2])
					} else if flag == '+' {
						lobbyId, _ := strconv.Atoi(args[1])
						if lobby, ok := l.lobbies[lobbyId]; ok {
							lobby.broadcast <- []byte(args[2])
						}
					}
				}
			}

			fmt.Println("Received msg from tcp server: ", string(req))
		}
	}
}

func (l *Leader) shutdown() {
	l.WSServer.Shutdown()
	l.TCPServer.Shutdown()
	close(l.TCPrequests)
	close(l.WSrequests)
}

// func (l *Leader) chooseNewLeader() {

// }

func idGenerator(beginnningId int) func() int {
	id := beginnningId
	return func() int {
		id++
		return id
	}
}

func (l *Leader) disconnectUser(username string) {
	user := l.Users[username]
	roomId := user.roomId
	if roomId == 1 {
		l.hub.unregister <- user
	} else {
		l.lobbies[roomId].unregister <- user
	}
	delete(l.Users, username)
}

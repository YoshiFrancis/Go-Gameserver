package leaderserver

import (
	"fmt"
	"strconv"

	"github.com/yoshifrancis/go-gameserver/src/messages"
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

type Leader struct {
	serverId     int                  // server id
	servers      map[int]bool         // server id, isExist
	clients      map[string]*User     // username, user pointer
	hub          *Hub                 // room title, room id
	lobbies      map[int]*Lobby       // lobby id, Lobby pointer
	applications map[int]*Application // app id, application pointer
	WSrequests   chan []byte          // channel of incoming requests from websocket server
	TCPrequests  chan []byte          // channel of incoming requests from tcp server
	WSServer     *wsserver.WSServer   // access to websocket server
	TCPServer    *tcpserver.TCPServer // access to tcp server
	idGen        func() int           // used to generate ids, get function from a closure
	isLeader     bool                 // indicate if leader or not
}

func NewLeader() *Leader {
	return &Leader{
		serverId:     1,
		servers:      make(map[int]bool),
		clients:      make(map[string]*User),
		hub:          NewHub(),
		lobbies:      make(map[int]*Lobby),
		applications: make(map[int]*Application),
		WSrequests:   make(chan []byte, 1024),
		TCPrequests:  make(chan []byte, 1024),
		idGen:        idGenerator(0),
		isLeader:     true,
	}
}

func (l *Leader) Run() {
	defer l.shutdown()

	for {
		select {
		case req := <-l.WSrequests:
			// if l.isLeader {
			// 	_, args := messages.Decode(req)
			// 	message := handleArgs(args)
			// 	l.TCPServer.Broadcast <- []byte(message)
			// } else {
			// 	l.TCPServer.Broadcast <- req
			// }
			fmt.Println("Received from WSServer", string(req))
			l.TCPServer.Broadcast <- req
		case req := <-l.TCPrequests:
			flag, args := messages.Decode(req)
			if l.isLeader {
				message := handleArgs(args)
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

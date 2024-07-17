package leaderserver

import (
	"fmt"

	"github.com/yoshifrancis/go-gameserver/src/messages"
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

type Leader struct {
	servers      map[int]bool         // server id, isExist
	clients      map[string]int       // username, server id
	hub          *Hub                 // room title, room id
	lobbies      map[int]*Lobby       // lobby id, Lobby pointer
	applications map[int]*Application // app id, application pointer
	requests     chan []byte          // channel of incomming requests from websocket server and external tcp servers
	WSServer     *wsserver.WSServer   // access to websocket server
	TCPServer    *tcpserver.TCPServer // access to tcp server
	idGen        func() int           // used to generate ids, get function from a closure
}

func NewLeader() *Leader {
	return &Leader{
		servers:      make(map[int]bool),
		clients:      make(map[string]int),
		hub:          NewHub(),
		lobbies:      make(map[int]*Lobby),
		applications: make(map[int]*Application),
		requests:     make(chan []byte, 1024),
		idGen:        idGenerator(),
	}
}

func (l *Leader) Run() {
	defer l.shutdown()

	for req := range l.requests {
		_, args := messages.Decode(req)
		message := handleArgs(args)
		fmt.Println(message)
	}
}

func (l *Leader) shutdown() {

}

func (l *Leader) chooseNewLeader() {

}

func idGenerator() func() int {
	id := 1
	return func() int {
		id++
		return id
	}
}

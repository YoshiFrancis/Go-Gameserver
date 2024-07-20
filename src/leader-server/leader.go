package leaderserver

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	fmt.Println("leader is running")
	defer l.shutdown()
	go l.hub.run()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case req := <-l.WSrequests:
			fmt.Println("received message from web socket server")
			broadcast_message := l.handleArgs(messages.Decode(req))
			l.TCPServer.Broadcast <- []byte(broadcast_message)
		case req := <-l.TCPrequests:
			fmt.Println("Received message from tcp server")
			flag, args := messages.Decode(req)
			l.handleArgs(flag, args)
		case sig := <-sigCh:
			fmt.Println("Received signal: ", sig)
			return
		}
	}
}

func idGenerator(beginnningId int) func() int {
	id := beginnningId
	return func() int {
		id++
		return id
	}
}

func (l *Leader) disconnectUser(username string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	user := l.Users[username]
	roomId := user.roomId
	if roomId == l.hub.hubId {
		l.hub.unregister <- user
	} else {
		lobby, ok := l.lobbies[roomId]
		if ok {
			lobby.unregister <- user
		}
	}
	delete(l.Users, username)
}

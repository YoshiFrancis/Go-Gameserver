package leaderserver

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/yoshifrancis/go-gameserver/src/messages"
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

type Leader struct {
	WSServer       *wsserver.WSServer   // access to websocket server
	TCPServer      *tcpserver.TCPServer // access to tcp server
	WSrequests     chan []byte          // channel of incoming requests from websocket server
	TCPrequests    chan []byte          // channel of incoming requests from tcp server
	NewConnections chan net.Conn        // channel of new connections to merge with
	isLeader       bool                 // indicate if leader or not
	idGen          func() int           // used to generate ids, get function from a closure
	hub            *Hub                 // room title, room id
	lobbies        map[int]*Lobby       // lobby id, Lobby pointer
	Users          map[string]*User     // username, user pointer
	mutex          sync.Mutex
}

func NewLeader() *Leader {
	return &Leader{
		Users:          make(map[string]*User),
		hub:            NewHub(0),
		lobbies:        make(map[int]*Lobby),
		WSrequests:     make(chan []byte, 1024),
		TCPrequests:    make(chan []byte, 1024),
		NewConnections: make(chan net.Conn, 10),
		idGen:          idGenerator(0),
		isLeader:       true,
		mutex:          sync.Mutex{},
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

func (l *Leader) sendAllData(conn net.Conn) {
	// send users
	// send hub -> [hubId]
	// send lobby -> [[lobbyId]]
	// user -> [username, serverId, roomId]
	hubId := l.hub.hubId
	lobbyIds := make([]int, 0)
	for lobbyId := range l.lobbies {
		lobbyIds = append(lobbyIds, lobbyId)
	}
	users := make([][]string, 0)
	for _, user := range l.Users {
		users = append(users, []string{user.username, strconv.Itoa(user.serverId), strconv.Itoa(user.roomId)})
	}
	starting_int := l.idGen()
	msg := messages.ServerMergeData(hubId, lobbyIds, users, starting_int)
	conn.Write([]byte(msg))

}

func (s *Leader) ConnectToServer(url string) bool {
	tcp := s.TCPServer
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting to server")
		return false
	}
	s.sendAllData(conn)
	var serverId int
	for {
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erorr while reading", err.Error())
			return false
		}
		_, args := messages.Decode(buffer)
		serverId, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("given invalid serverid")
			return false
		}
		if strings.ToLower(args[0]) == "merge" {
			s.TCPrequests <- buffer
			break
		}
	}

	tcp.Register <- tcpserver.NewExternalTCPServer(tcp, conn, url, serverId)
	return true
}

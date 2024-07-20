package tcpserver

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

type TCPServer struct {
	Servers    map[int]*ExtenalTCPServer
	Broadcast  chan []byte
	requests   chan []byte
	Register   chan *ExtenalTCPServer
	unregister chan *ExtenalTCPServer
	serverId   int
	idGen      func() int
}

func NewTCPServer(requests chan []byte) *TCPServer {
	return &TCPServer{
		Servers:    make(map[int]*ExtenalTCPServer),
		Broadcast:  make(chan []byte, 1024),
		requests:   requests,
		Register:   make(chan *ExtenalTCPServer, 10),
		unregister: make(chan *ExtenalTCPServer, 10),
		serverId:   0,
		idGen:      idGeneratorInit(1),
	}
}

func (s *TCPServer) Listen(url string) {
	listener, err := net.Listen("tcp", url)
	if err != nil {
		fmt.Println("Error getting listener socket for tcp server", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Beginning to listen for other servers at the url", url)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}

		fmt.Println("Connection from ", conn.LocalAddr().String())
		s.Register <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), s.idGen())
	}
}

func (s *TCPServer) Run() {
	for {
		select {
		case message := <-s.Broadcast:
			for _, server := range s.Servers {
				server.Send <- message
			}
		case c := <-s.Register:
			fmt.Println("Server registered!")
			s.Servers[c.serverId] = c
			go c.run()
			s.Broadcast <- []byte(messages.ServerCreation(c.serverId, c.Url))
		case c := <-s.unregister:
			delete(s.Servers, c.serverId)
			fmt.Println("Server unregistered")
		}
	}
}

func (s *TCPServer) Shutdown() {
	fmt.Println("tcpserver shutting down")
	close(s.Broadcast)
	close(s.requests)
	close(s.Register)
	close(s.unregister)
}

func (s *TCPServer) ConnectToServer(url string) bool {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting to server")
		return false
	}

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
		if strings.ToLower(args[0]) == "creation" {
			break
		}
	}

	s.Register <- NewExternalTCPServer(s, conn, url, serverId)
	s.idGen = idGeneratorInit(serverId)
	return true
}

func idGeneratorInit(start int) func() int {
	i := start
	return func() int {
		i++
		return i
	}
}

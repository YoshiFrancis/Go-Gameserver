package tcpserver

import (
	"fmt"
	"net"
)

type TCPServer struct {
	Servers    map[int]*ExtenalTCPServer
	Broadcast  chan []byte
	requests   chan []byte
	register   chan *ExtenalTCPServer
	unregister chan *ExtenalTCPServer
	serverId   int
}

func NewTCPServer(requests chan []byte) *TCPServer {
	return &TCPServer{
		Servers:    make(map[int]*ExtenalTCPServer),
		Broadcast:  make(chan []byte, 1024),
		requests:   requests,
		register:   make(chan *ExtenalTCPServer, 10),
		unregister: make(chan *ExtenalTCPServer, 10),
		serverId:   -1,
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
		s.register <- NewExternalTCPServer(s, conn, conn.LocalAddr().String())
	}
}

func (s *TCPServer) Run() {
	defer s.Shutdown()
	for {
		select {
		case message := <-s.Broadcast:
			for _, server := range s.Servers {
				server.Send <- message
			}
		// case message := <-s.requests:
		// 	_, args := messages.Decode(message)
		// 	fmt.Println("Received ", args)
		case c := <-s.register:
			fmt.Println("Server registered!")
			s.Servers[c.serverId] = c
			go c.run()
		case c := <-s.unregister:
			delete(s.Servers, c.serverId)
			fmt.Println("Server unregistered")
		}
	}
}

func (s *TCPServer) Shutdown() {
	close(s.Broadcast)
	close(s.requests)
	close(s.register)
	close(s.unregister)
}

func (s *TCPServer) ConnectToServer(url string) bool {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting to server")
		return false
	}
	s.register <- NewExternalTCPServer(s, conn, url)
	return true
}

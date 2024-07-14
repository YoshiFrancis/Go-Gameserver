package tcpserver

import (
	"fmt"
	"net"
)

type TCPServer struct {
	servers    map[*ExtenalTCPServer]bool
	broadcast  chan byte
	read       chan byte
	register   chan *ExtenalTCPServer
	unregister chan *ExtenalTCPServer
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		servers:    make(map[*ExtenalTCPServer]bool),
		broadcast:  make(chan byte),
		read:       make(chan byte),
		register:   make(chan *ExtenalTCPServer),
		unregister: make(chan *ExtenalTCPServer),
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

		s.register <- NewExternalTCPServer(s, conn, url)
	}
}

func (s *TCPServer) Run() {
	for {
		select {
		case message := <-s.broadcast:
			fmt.Println("Ready to send ", message)
		case message := <-s.read:
			fmt.Println("Received ", message)
		case c := <-s.register:
			fmt.Println("Server registered!")
			s.servers[c] = true
			c.run()
		case c := <-s.unregister:
			delete(s.servers, c)
			fmt.Println("Server unregistered")
		}
	}
}

func (s *TCPServer) Close() {
	close(s.broadcast)
	close(s.read)
	close(s.register)
	close(s.unregister)
}

func (s *TCPServer) ConnectToServer(url string) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting to server")
		return
	}

	s.register <- NewExternalTCPServer(s, conn, url)
}

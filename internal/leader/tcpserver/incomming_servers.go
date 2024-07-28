package tcpserver

import (
	"fmt"
	"net"
)

// TODO
// adjust file to accomodate changes to TCPServer class
// also make struct to encapsulate how to swend messages to leaders and followers
// so creation of LeaderMessage and FollowerMessage (or perhaps just a Message)

type ExtenalTCPServer struct {
	main_server *TCPServer
	conn        net.Conn
	Send        chan []byte
	Shutdown    chan bool
	Url         string
	serverId    int
	class       string
}

func NewExternalTCPServer(main_server *TCPServer, conn net.Conn, url string, serverId int, class string) *ExtenalTCPServer {
	return &ExtenalTCPServer{
		main_server: main_server,
		conn:        conn,
		Send:        make(chan []byte, 1024),
		Shutdown:    make(chan bool, 1),
		Url:         url,
		serverId:    serverId,
		class:       class,
	}
}

func (s *ExtenalTCPServer) run() {
	defer s.close()
	go s.read()
	for {
		select {
		case send := <-s.Send:
			fmt.Println("external tcpserver sending message")
			_, err := s.conn.Write(send)
			if err != nil {
				fmt.Println("Error sending to server at: ", s.Url)
			}
		case <-s.Shutdown:
			break
		}
	}
}

func (s *ExtenalTCPServer) read() {
	for {
		buffer := make([]byte, 1024)
		_, err := s.conn.Read(buffer)
		if err != nil {
			fmt.Println("Erorr while reading", err.Error())
			s.Shutdown <- true
			return
		}

		if s.class == "F" {
			s.main_server.fRequests <- buffer
		} else {
			s.main_server.fRequests <- buffer
		}

	}
}

func (s *ExtenalTCPServer) close() {
	close(s.Send)
	close(s.Shutdown)
	if s.class == "F" {
		s.main_server.fRegistry <- s
	} else {
		s.main_server.lRegistry <- s
	}
	s.conn.Close()
	close(s.Send)
}

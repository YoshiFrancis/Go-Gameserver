package tcpserver

import (
	"fmt"
	"net"
)

type ExtenalTCPServer struct {
	main_server *TCPServer
	conn        net.Conn
	Send        chan []byte
	Shutdown    chan bool
	Url         string
}

func NewExternalTCPServer(main_server *TCPServer, conn net.Conn, url string) *ExtenalTCPServer {
	return &ExtenalTCPServer{
		main_server: main_server,
		conn:        conn,
		Send:        make(chan []byte, 1024),
		Shutdown:    make(chan bool, 1),
		Url:         url,
	}
}

func (s *ExtenalTCPServer) run() {
	defer s.close()
	go s.read()
	for {
		select {
		case send := <-s.Send:
			fmt.Println("tcpserver sending message!")
			_, err := s.conn.Write(send)
			if err != nil {
				fmt.Println("Error sending to server at: ", s.Url)
			}
		case <-s.Shutdown:
			return
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
		s.main_server.requests <- buffer
	}
}

func (s *ExtenalTCPServer) close() {
	s.main_server.unregister <- s
	s.conn.Close()
	close(s.Send)
}

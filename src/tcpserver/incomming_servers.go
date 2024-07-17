package tcpserver

import (
	"fmt"
	"net"
)

type ExtenalTCPServer struct {
	main_server *TCPServer
	conn        net.Conn
	send        chan []byte
	shutdown    chan bool
	url         string
}

func NewExternalTCPServer(main_server *TCPServer, conn net.Conn, url string) *ExtenalTCPServer {
	return &ExtenalTCPServer{
		main_server: main_server,
		conn:        conn,
		send:        make(chan []byte, 1024),
		shutdown:    make(chan bool, 1),
		url:         url,
	}
}

func (s *ExtenalTCPServer) run() {
	defer s.close()
	go s.read()
	for {
		select {
		case send := <-s.send:
			fmt.Println("tcpserver sending message!")
			_, err := s.conn.Write(send)
			if err != nil {
				fmt.Println("Error sending to server at: ", s.url)
			}
		case <-s.shutdown:
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
			s.shutdown <- true
			return
		}
		s.main_server.requests <- buffer
	}
}

func (s *ExtenalTCPServer) close() {
	s.main_server.unregister <- s
	s.conn.Close()
	close(s.send)
}

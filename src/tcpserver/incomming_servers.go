package tcpserver

import (
	"fmt"
	"net"
)

type ExtenalTCPServer struct {
	main_server *TCPServer
	conn        net.Conn
	send        chan []byte
	url         string
}

func NewExternalTCPServer(main_server *TCPServer, conn net.Conn, url string) *ExtenalTCPServer {
	return &ExtenalTCPServer{
		main_server: main_server,
		conn:        conn,
		send:        make(chan []byte, 1024),
		url:         url,
	}
}

func (s *ExtenalTCPServer) run() {
	defer s.close()
	for send := range s.send {
		_, err := s.conn.Write(send)
		if err != nil {
			fmt.Println("Error sending to server at: ", s.url)
		}
	}
}

func (s *ExtenalTCPServer) close() {
	s.main_server.unregister <- s
	s.conn.Close()
	close(s.send)
}

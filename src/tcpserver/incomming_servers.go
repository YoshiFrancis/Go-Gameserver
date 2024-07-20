package tcpserver

import (
	"fmt"
	"net"
	"strconv"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

type ExtenalTCPServer struct {
	main_server *TCPServer
	conn        net.Conn
	Send        chan []byte
	Shutdown    chan bool
	Url         string
	serverId    int
}

func NewExternalTCPServer(main_server *TCPServer, conn net.Conn, url string, serverId int) *ExtenalTCPServer {
	return &ExtenalTCPServer{
		main_server: main_server,
		conn:        conn,
		Send:        make(chan []byte, 1024),
		Shutdown:    make(chan bool, 1),
		Url:         url,
		serverId:    serverId,
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
		fmt.Println("External tcp server received message")
		flag, msg := messages.Decode(buffer)
		if flag == '-' {
			if msg[0] == "serverid" {
				serverId, err := strconv.Atoi(msg[1])
				if err != nil {
					fmt.Println("given invalid server id")
					continue
				}
				s.serverId = serverId
			} else if msg[0] == "accept" {
				serverId, err := strconv.Atoi(msg[1])
				if err != nil {
					fmt.Println("given invalid server id")
					continue
				}
				url := msg[2]
				s.main_server.AcceptConnectedServer(serverId, url)
			}
		} else {
			s.main_server.requests <- buffer
		}
	}
}

func (s *ExtenalTCPServer) close() {
	s.main_server.unregister <- s
	s.conn.Close()
	close(s.Send)
}

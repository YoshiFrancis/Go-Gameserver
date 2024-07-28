package tcpserver

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
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
		flag, args := messages.Decode(buffer)
		args[0] = strings.ToLower(args[0])
		if flag == '-' {
			switch args[0] {
			case "serverid":
				serverId, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("given invalid server id")
					continue
				}
				s.serverId = serverId
			case "accept":
				serverId, err := strconv.Atoi(args[1])
				if err != nil {
					fmt.Println("given invalid server id")
					continue
				}
				url := args[2]
				s.main_server.AcceptConnectedServer(serverId, url)
			case "shutdown":
				s.main_server.unregister <- s

				// the server that was originally connected now must broadcast to all other servers rhat there is a neew server
				// i have to come up with new key word to signal that the new server has already been accepted by one of the nodes in the group already
				// the new servers will connect with the already connected node
				// this original node that accepted has the send all data about the servers to tje new node

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

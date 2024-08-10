package tcpserver

import (
	"fmt"
	"net"

	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

// TODO
// adjust file to accomodate changes to TCPServer class
// also make struct to encapsulate how to swend messages to leaders and followers
// so creation of LeaderMessage and FollowerMessage (or perhaps just a Message)

type Request struct {
	flag    byte
	command string
	arg     string
	sender  string
	server  *ExtenalTCPServer
}

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
			if s.Shutdown != nil {
				return
			}
		}
		fReq := messages.FReqDecode(buffer)

		new_req := Request{
			flag:    fReq.Flag,
			command: fReq.Command,
			arg:     fReq.Arg,
			sender:  fReq.Sender,
			server:  s,
		}

		if new_req.flag == '!' && new_req.arg == "ping" {
			s.Send <- []byte(messages.Pong())
		} else if s.class == "F" {
			s.main_server.fRequests <- new_req
		} else {
			s.main_server.lRequests <- buffer
		}

	}
}

func (s *ExtenalTCPServer) close() {

	fmt.Println("Follower server is closing down in leader")
	close(s.Send)
	close(s.Shutdown)
	if s.class == "F" {
		s.main_server.fRegistry <- s
	} else {
		s.main_server.lRegistry <- s
	}
	s.conn.Close()
}

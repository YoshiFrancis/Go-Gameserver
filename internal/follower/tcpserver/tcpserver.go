package tcpserver

import (
	"fmt"
	"net"

	"github.com/yoshifrancis/go-gameserver/internal/follower/wsserver"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

type TCPServer struct {
	Lto        chan []byte
	Lfrom      chan wsserver.LeaderRequest
	WSfrom     chan []byte
	WSto       chan wsserver.LeaderRequest
	done       chan bool
	leaderConn net.Conn
	ServerId   string
}

func NewTCPServer(conn net.Conn, serverId string, done chan bool) *TCPServer {
	return &TCPServer{
		Lto:   make(chan []byte),
		Lfrom: make(chan wsserver.LeaderRequest, 8),
		done:  done,

		leaderConn: conn,
		ServerId:   serverId,
	}
}

func ConnectToLeader(leaderIp string, done chan bool) *TCPServer {
	conn, err := net.Dial("tcp", leaderIp)
	if err != nil {
		fmt.Println("Error connecting to leader: ", err.Error())
		return nil
	}
	return NewTCPServer(conn, "", done) // hardcoded server id
}

func (tcp *TCPServer) Run() {
	defer tcp.Shutdown()
	go tcp.leaderRead()
	for {
		select {
		case lReq := <-tcp.Lfrom: // -------- these two lines somehow cause a bug (cannot switch from username.html to hub.html)
			tcp.WSto <- lReq
		case from := <-tcp.WSfrom:
			_, err := tcp.leaderConn.Write(from)
			if err != nil {
				fmt.Println("Error writing to leader: ", err)
				return
			}
		}
	}
}

func (s *TCPServer) Shutdown() {
	fmt.Println("tcpserver shutting down")
	close(s.Lto)
	close(s.Lfrom)
	close(s.WSfrom)
	close(s.WSto)
	s.leaderConn.Close()
	if s.done != nil {
		s.done <- true
	}
}

func (s *TCPServer) leaderRead() {
	for {
		buffer := make([]byte, 1024)
		_, err := s.leaderConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from leader: ", err)
			return
		}

		go func() {
			lReq := messages.LReqDecode(buffer)
			s.Lfrom <- wsserver.NewLeaderRequest(lReq.Command, lReq.Arg, lReq.Receivers)
		}()
	}
}

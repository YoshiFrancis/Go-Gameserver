package tcpserver

import (
	"bytes"
	"fmt"
	"net"
)

type TCPServer struct {
	Lto        chan []byte
	Lfrom      chan []byte
	WSfrom     chan []byte
	WSto       chan []byte
	leaderConn net.Conn
	ServerId   string
}

func NewTCPServer(conn net.Conn, serverId string) *TCPServer {
	return &TCPServer{
		Lto:   make(chan []byte),
		Lfrom: make(chan []byte),

		leaderConn: conn,
		ServerId:   serverId,
	}
}

func ConnectToLeader(leaderIp string) *TCPServer {
	conn, err := net.Dial("tcp", leaderIp)
	if err != nil {
		fmt.Println("Error connecting to leader: ", err.Error())
		return nil
	}
	return NewTCPServer(conn, "") // hardcoded server id
}

func (tcp *TCPServer) Run() {
	defer tcp.Shutdown()
	go tcp.leaderRead()
	for {
		select {
		// 	case from := <-tcp.Lfrom: -------- these two lines somehow cause a bug (cannot switch from username.html to hub.html)
		// 		tcp.WSto <- from
		case from := <-tcp.WSfrom:
			fmt.Println("received message from WS: ", string(from))
			fmt.Println("WRITING TO LEADER")
			_, err := tcp.leaderConn.Write(from)
			if err != nil {
				fmt.Println("Error writing to leader: ", err)
				return
			}
			fmt.Println("TCP Giving to send to leader")
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
}

func (s *TCPServer) leaderRead() {
	for {
		var buffer bytes.Buffer
		_, err := s.leaderConn.Read(buffer.Bytes())
		if err != nil {
			fmt.Println("Error reading from leader: ", err)
			return
		}
		s.Lfrom <- buffer.Bytes()
	}
}

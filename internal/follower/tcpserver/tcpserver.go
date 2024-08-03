package tcpserver

import (
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
		case from := <-tcp.Lfrom: // -------- these two lines somehow cause a bug (cannot switch from username.html to hub.html)
			tcp.WSto <- from
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
}

func (s *TCPServer) leaderRead() {
	for {
		buffer := make([]byte, 1024)
		n, err := s.leaderConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from leader: ", err)
			return
		}
		s.Lfrom <- buffer[:n]
	}
}

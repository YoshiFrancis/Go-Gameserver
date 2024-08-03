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
		case to := <-tcp.Lto:
			_, err := tcp.leaderConn.Write(to)
			if err != nil {
				fmt.Println("Error writing to leader: ", err)
			}
		case from := <-tcp.Lfrom:
			tcp.WSto <- from
		case from := <-tcp.WSfrom:
			tcp.Lto <- from
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

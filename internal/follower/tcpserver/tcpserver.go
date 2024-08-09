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
	done       chan bool
	leaderConn net.Conn
	ServerId   string
}

func NewTCPServer(conn net.Conn, serverId string, done chan bool) *TCPServer {
	return &TCPServer{
		Lto:   make(chan []byte),
		Lfrom: make(chan []byte),
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
		case from := <-tcp.Lfrom: // -------- these two lines somehow cause a bug (cannot switch from username.html to hub.html)
			fmt.Println("Got message from leader!", from)
			tcp.WSto <- from
			fmt.Println("Sent to WS")
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
		n, err := s.leaderConn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from leader: ", err)
			return
		}
		s.Lfrom <- buffer[:n]
	}
}

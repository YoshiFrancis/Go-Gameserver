package tcpserver

import (
	"fmt"
	"net"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

type TCPServer struct {
	Broadcast chan []byte
	Lfrom     chan []byte
	WSfrom    chan []byte
	WSto      chan []byte
	conn      net.Conn
	ServerId  string
}

func NewTCPServer(conn net.Conn, serverId string) *TCPServer {
	return &TCPServer{
		Broadcast: make(chan []byte, 1024),
		Lfrom:     make(chan []byte, 1024),
		conn:      conn,
		ServerId:  serverId,
	}
}

func ConnectToLeader(leaderIp string) *TCPServer {
	conn, err := net.Dial("tcp", leaderIp)
	if err != nil {
		fmt.Println("Error connecting to leader: ", err.Error())
		return nil
	}
	severId := getServerId(conn)
	return NewTCPServer(conn, severId)
}

func (tcp *TCPServer) Run() {

}

func getServerId(conn net.Conn) string {
	for {
		message := make([]byte, 50)
		_, err := conn.Read(message)
		if err != nil {
			fmt.Println("Client is going to stop reading!")
			conn.Close()
			return ""
		}
		_, args := messages.Decode(message)
		if strings.ToLower(args[0]) == "serverid" {
			return args[1]
		}
	}
}

func (s *TCPServer) Shutdown() {
	fmt.Println("tcpserver shutting down")
	close(s.Broadcast)
	close(s.Lfrom)
	close(s.WSfrom)
	close(s.WSto)
}

package tcpserver

import (
	"fmt"
	"net"

	"github.com/yoshifrancis/go-gameserver/src/messages"
)

type TCPServer struct {
	lServers  map[string]*ExtenalTCPServer // name of server, tcp connection container
	fServers  map[int]*ExtenalTCPServer    // id of follower server, tcp connnection container
	lRegistry chan *ExtenalTCPServer
	fRegistry chan *ExtenalTCPServer
	lRequests chan []byte
	fRequests chan []byte
	serverId  int
	idGen     func() int
}

func NewTCPServer(requests chan []byte) *TCPServer {
	return &TCPServer{
		lRegistry: make(chan *ExtenalTCPServer, 5),
		fRegistry: make(chan *ExtenalTCPServer, 5),
		lRequests: make(chan []byte, 1024),
		fRequests: make(chan []byte, 1024),
		serverId:  0,
	}
}

// TODO
// listen for leader and follower and able to distingusih between the two
// do follower first
// implement tests somehow
func (s *TCPServer) Listen(url string) {
	listener, err := net.Listen("tcp", url)
	if err != nil {
		fmt.Println("Error getting listener socket for tcp server", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Beginning to listen for other servers at the url", url)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}

		fmt.Println("Connection from ", conn.LocalAddr().String())
		s.register <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), s.idGen())
		s.newConnections <- conn
	}
}

// TODO
// remove all followers
// tell all leaders (for later)
func (s *TCPServer) Shutdown() {
	fmt.Println("tcpserver shutting down")

}

// TODO
// listen for follower requests and registry
// move on to leader requests and registry
func (s *TCPServer) Run() {

}

// TODO
// able to distinguish between follower and leader
// perhaps two differnet ports for follower and leaders would work
func (s *TCPServer) AcceptConnectedServer(serverId int, url string) bool {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Println("Error connecting to server")
		return false
	}

	s.fRegistry <- NewExternalTCPServer(s, conn, url, serverId)
	msg := messages.ServerTellServerId(serverId)
	conn.Write([]byte(msg))
	return true
}

func (s *TCPServer) ServerId() int {
	return s.serverId
}

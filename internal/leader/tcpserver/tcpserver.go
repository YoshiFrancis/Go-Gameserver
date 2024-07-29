package tcpserver

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type TCPServer struct {
	lServers  map[string]*ExtenalTCPServer // name of server, tcp connection container
	fServers  map[int]*ExtenalTCPServer    // id of follower server, tcp connnection container
	lRegistry chan *ExtenalTCPServer
	fRegistry chan *ExtenalTCPServer
	lRequests chan []byte
	fRequests chan []byte
	url       string
	idGen     func() int
	mux       sync.Mutex
}

func NewTCPServer() *TCPServer {
	return &TCPServer{
		lServers:  make(map[string]*ExtenalTCPServer),
		fServers:  make(map[int]*ExtenalTCPServer),
		lRegistry: make(chan *ExtenalTCPServer, 5),
		fRegistry: make(chan *ExtenalTCPServer, 5),
		lRequests: make(chan []byte, 1024),
		fRequests: make(chan []byte, 1024),

		mux: sync.Mutex{},
	}
}

func (s *TCPServer) Run() {
	defer s.Shutdown()
	for {
		leaderUrl := promptText("Input port to listen for other servers or -1 if no server-> :<PORT>")
		if leaderUrl == ":-1" {
			break
		}
		leaderListener, err := net.Listen("tcp", leaderUrl)
		if err != nil {
			fmt.Println("Error setting up listener for tcp server", err.Error())
			continue
		}
		go s.listenForLeaders(leaderListener)
		break
	}

	for {
		followerUrl := promptText("Input port to listen for follower -> :<PORT>")
		followerListener, err := net.Listen("tcp", followerUrl)
		if err != nil {
			fmt.Println("Error setting up listener for tcp server", err.Error())
			continue
		}

		go s.listenForFollower(followerListener)
		break
	}

	for {
		select {
		case l := <-s.lRegistry:
			s.mux.Lock()

			if server, ok := s.lServers[l.conn.LocalAddr().String()]; ok {
				server.Shutdown <- true
				delete(s.lServers, l.conn.LocalAddr().String())
			} else {
				s.lServers[l.conn.LocalAddr().String()] = l
			}

			s.mux.Unlock()
			// tell followers and leaders I guess
		case f := <-s.fRegistry:
			s.mux.Lock()

			if server, ok := s.fServers[f.serverId]; ok {
				server.Shutdown <- true
				delete(s.fServers, f.serverId)
			} else {
				s.fServers[f.serverId] = f
			}
			go f.run()

			s.mux.Unlock()
			// do not need to tell anyone about them
		case lReq := <-s.lRequests:
			fmt.Println(lReq)
			// handle leader requests
		case fReq := <-s.fRequests:
			fmt.Println(fReq)
			// handle request from follower
		}
	}
}

func (s *TCPServer) Input() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadBytes('\n')
		text = []byte(strings.Replace(string(text), "\n", "", -1))
		text_string := string(text)
		fmt.Println("handle: ", text_string) // TODO
	}
}

func promptText(prompt string) string {
	fmt.Println(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadBytes('\n')
	text = []byte(strings.Replace(string(text), "\n", "", -1))
	return string(text)
}

func (s *TCPServer) listenForFollower(listener net.Listener) bool {

	defer listener.Close()
	fmt.Println("leader beginning to listen for followers at: ", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}
		s.fRegistry <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), s.idGen(), "F")
	}
}

func (s *TCPServer) listenForLeaders(listener net.Listener) {
	defer listener.Close()
	s.url = listener.Addr().String()
	fmt.Println("leader beginning to listen for other leaders at: ", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}
		s.lRegistry <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), s.idGen(), "L")
	}
}

// TODO
// remove all followers
// tell all leaders (for later)
func (s *TCPServer) Shutdown() {
	fmt.Println("tcpserver shutting down")
	for _, l := range s.lServers {
		l.Shutdown <- true
	}

	for _, f := range s.fServers {
		f.Shutdown <- true
	}

	close(s.fRequests)
	close(s.lRequests)
	close(s.lRegistry)
	close(s.fRegistry)
}

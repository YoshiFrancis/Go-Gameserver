package tcpserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
	"github.com/yoshifrancis/go-gameserver/internal/leader/rooms"
)

type TCPServer struct {
	lServers    map[string]*ExtenalTCPServer // name of server, tcp connection container to leader
	fServers    map[string]*ExtenalTCPServer // id of follower server, tcp connnection container to follower
	aServers    map[string]*ExtenalTCPServer // name of application, tcp connection connection to application
	lRegistry   chan *ExtenalTCPServer
	fRegistry   chan *ExtenalTCPServer
	aRegistry   chan *ExtenalTCPServer
	lRequests   chan []byte
	fRequests   chan FollowerRequest
	aRequest    chan ApplicationRequest
	idGen       func() int
	url         string
	userStorage *containers.Storage[string, *rooms.User] // username, User struct
	roomStorage *containers.Storage[string, rooms.Room]  // roomId, room
	done        chan bool
	hub         *rooms.Hub
	mux         sync.Mutex
}

func NewTCPServer(done chan bool) *TCPServer {
	new_server := &TCPServer{
		lServers:    make(map[string]*ExtenalTCPServer),
		fServers:    make(map[string]*ExtenalTCPServer),
		aServers:    make(map[string]*ExtenalTCPServer),
		lRegistry:   make(chan *ExtenalTCPServer),
		fRegistry:   make(chan *ExtenalTCPServer),
		aRegistry:   make(chan *ExtenalTCPServer),
		lRequests:   make(chan []byte),
		fRequests:   make(chan FollowerRequest),
		aRequest:    make(chan ApplicationRequest),
		idGen:       idGen(),
		userStorage: containers.NewStorage[string, *rooms.User](),
		roomStorage: containers.NewStorage[string, rooms.Room](),
		done:        done,
		hub:         rooms.NewHub(1),
		mux:         sync.Mutex{},
	}
	new_server.roomStorage.Set("hub", new_server.hub)
	return new_server
}

func (s *TCPServer) Run(tcpFPort, tcpAPort string) {
	defer s.Shutdown()

	followerListener, err := net.Listen("tcp", tcpFPort)
	if err != nil {
		fmt.Println("Error setting up listener for tcp server", err.Error())
		panic(err)
	}

	applicationListener, err := net.Listen("tcp", tcpAPort)
	if err != nil {
		fmt.Println("Error setting up listener for tcp server", err.Error())
		panic(err)
	}

	go s.listenForFollower(followerListener)
	go s.listenForApplications(applicationListener)

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
				delete(s.fServers, server.serverId)
			} else {
				s.fServers[f.serverId] = f
				go f.run()
			}

			s.mux.Unlock()
			// do not need to tell anyone about them
		case a := <-s.aRegistry:
			s.mux.Lock()

			if server, ok := s.aServers[a.serverId]; ok {
				delete(s.aServers, server.serverId)
			} else {
				s.aServers[a.serverId] = a
				go a.run()
			}

			s.mux.Unlock()
		case lReq := <-s.lRequests:
			fmt.Println("Request from another leader: ", string(lReq))
			// handle leader requests
		case fReq := <-s.fRequests:
			s.handleFollowerRequest(fReq)
		case aReq := <-s.aRequest:
			s.handleApplicationRequest(aReq)
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

func (s *TCPServer) listenForFollower(listener net.Listener) bool {
	defer listener.Close()
	s.url = listener.Addr().String()
	fmt.Println("leader beginning to listen for followers at: ", s.url)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}
		fmt.Println("a follower connected!")
		s.fRegistry <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), string(s.idGen()), "F")
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
		s.lRegistry <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), string(s.idGen()), "L")
	}
}

func (s *TCPServer) listenForApplications(listener net.Listener) {
	defer listener.Close()
	fmt.Println("leader beginning to listen for applications at: ", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting new connection!")
			continue
		}

		log.Println("New application connected")

		// get username
		nameBuffer := make([]byte, 128)
		n, err := conn.Read(nameBuffer)
		if err != nil {
			log.Println("Error reading from application!")
			conn.Close()
			continue
		}

		s.aRegistry <- NewExternalTCPServer(s, conn, conn.LocalAddr().String(), string(nameBuffer[:n]), "A")

	}
}

func idGen() func() int {
	i := 1
	return func() int {
		i++
		return i
	}
}

// TODO
// remove all followers
// tell all leaders (for later)
func (s *TCPServer) Shutdown() {
	fmt.Println("leader tcpserver shutting down")
	for _, l := range s.lServers {
		l.Shutdown <- true
	}

	for _, f := range s.fServers {
		f.Shutdown <- true
	}

	waiting_for_shutdowns := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-waiting_for_shutdowns.C:
			if len(s.lServers) == 0 && len(s.fServers) == 0 {
				fmt.Println("leader shutting down reugquests and registries")
				close(s.fRequests)
				close(s.lRequests)
				close(s.lRegistry)
				close(s.fRegistry)

				if s.done != nil {
					s.done <- true
				}
				return
			}
		}
	}

}

func (s *TCPServer) fbroadcast(message string) {
	fmt.Println("Broadcasting: ", message)
	for _, server := range s.fServers {
		server.Send <- []byte(message)
	}
}

func (s *TCPServer) abroadcast(appName, message string) {
	app, ok := s.aServers[appName]
	if ok {
		app.Send <- []byte(message)
	}
}

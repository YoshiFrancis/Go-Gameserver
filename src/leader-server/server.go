package leaderserver

import (
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
)

type Server struct {
	serverId  int
	tcpserver *tcpserver.ExtenalTCPServer
}

func NewServer(serverId int, tcpserver *tcpserver.ExtenalTCPServer) *Server {
	return &Server{
		serverId:  serverId,
		tcpserver: tcpserver,
	}
}

func (server *Server) send(msg []byte) {
	server.tcpserver.Send <- msg
}

func (server *Server) shutdown() {
	server.tcpserver.Shutdown <- true
}

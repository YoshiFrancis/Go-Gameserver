package leader

import (
	"github.com/yoshifrancis/go-gameserver/internal/leader/tcpserver"
)

func Leader_init(tcpPort string, done chan bool) *tcpserver.TCPServer {
	leader := tcpserver.NewTCPServer(done)
	go leader.Run(tcpPort)
	return leader
}

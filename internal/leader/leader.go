package leader

import (
	"github.com/yoshifrancis/go-gameserver/internal/leader/tcpserver"
)

func Leader_init(tcpFPort, tcpAPort string, done chan bool) *tcpserver.TCPServer {
	leader := tcpserver.NewTCPServer(done)
	go leader.Run(tcpFPort, tcpAPort)
	return leader
}

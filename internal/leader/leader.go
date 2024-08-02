package leader

import (
	"github.com/yoshifrancis/go-gameserver/internal/leader/tcpserver"
)

func Leader_init(tcpPort string) {
	leader := tcpserver.NewTCPServer()
	go leader.Run(tcpPort)
}

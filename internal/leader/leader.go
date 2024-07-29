package leader

import (
	"github.com/yoshifrancis/go-gameserver/internal/leader/tcpserver"
)

func Leader_init(tcpPort int) {
	leader := tcpserver.NewTCPServer()
	go leader.Run()
}

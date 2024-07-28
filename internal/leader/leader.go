package leader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yoshifrancis/go-gameserver/internal/leader/tcpserver"
)

func Leader_init(tcpPort int) {
	leader := tcpserver.NewTCPServer()
	go leader.Run()
	input(leader)
}

func input(leader *tcpserver.TCPServer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadBytes('\n')
		text = []byte(strings.Replace(string(text), "\n", "", -1))
		text_string := string(text)
		fmt.Println("handle: ", text_string)
	}
}

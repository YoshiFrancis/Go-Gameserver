package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	leaderserver "github.com/yoshifrancis/go-gameserver/src/leader-server"
	"github.com/yoshifrancis/go-gameserver/src/messages"
	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

func main() {

	if len(os.Args) != 3 {
		log.Fatal("Usgae: go run main.go <WS PORT> < TCPPORT>")
	}

	wsPort := ":" + os.Args[1]
	tcpPort := ":" + os.Args[2]

	leader := leaderserver.NewLeader()

	wsserver := wsserver.NewWSServer(leader.WSrequests)
	tcpserver := tcpserver.NewTCPServer(leader.TCPrequests)
	leader.WSServer = wsserver
	leader.TCPServer = tcpserver
	go leader.Run()
	go wsserver.Run()
	go tcpserver.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsserver.Serve(w, r)
	})

	go input(tcpserver)
	go tcpserver.Listen(tcpPort)
	log.Fatal(http.ListenAndServe(wsPort, nil))
}

func input(tcpserver *tcpserver.TCPServer) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadBytes('\n')
		text = []byte(strings.Replace(string(text), "\n", "", -1))
		text_string := string(text)
		if text_string[0] == '/' {
			fmt.Println("Attempting to connect to ", text_string[:1])
			go func() {
				connected := tcpserver.ConnectToServer(text_string[1:])
				if connected {
					fmt.Println("Successfully connected to server: ", text_string[1:])
				} else {
					fmt.Println("Failed to connect to server: ", text_string[1:])
				}
			}()
		} else if text_string[0] == '\\' {
			msg := messages.HubBroadcast("server", 1, text_string[1:])
			tcpserver.Broadcast <- []byte(msg)
		} else {
			fmt.Println("Invalid command")
		}
	}
}

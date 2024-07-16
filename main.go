package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yoshifrancis/go-gameserver/src/tcpserver"
	"github.com/yoshifrancis/go-gameserver/src/wsserver"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usgae: go run main.go <PORT>")
	}

	listeningUrl := ":" + os.Args[1]

	wsserver := wsserver.NewServer()
	tcpserver := tcpserver.NewTCPServer()
	wsserver.TCPRead = tcpserver.WSSend
	tcpserver.WSRead = wsserver.TCPSend
	go wsserver.Run()
	go tcpserver.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		wsserver.Serve(w, r)
	})

	go input(tcpserver, wsserver)
	go tcpserver.Listen(listeningUrl)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func input(tcpserver *tcpserver.TCPServer, wsserver *wsserver.Server) {
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
			tcpserver.WSRead <- []byte(text_string[1:])
			wsserver.TCPRead <- []byte(text_string[1:])
		}
	}
}

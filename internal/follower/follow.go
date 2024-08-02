package follower

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yoshifrancis/go-gameserver/internal/follower/tcpserver"
	"github.com/yoshifrancis/go-gameserver/internal/follower/wsserver"
)

func Follower_init(wsPort, leaderIp string) {
	tcp := tcpserver.ConnectToLeader(leaderIp)
	if tcp == nil {
		fmt.Println("Error connecting to leader")
		return
	}
	fmt.Println("Connected to leader")

	ws := wsserver.NewWSServer()
	link_1 := make(chan []byte, 1024)
	link_2 := make(chan []byte, 1024)
	tcp.WSfrom = link_1
	ws.TCPto = link_1
	tcp.WSto = link_2
	ws.TCPfrom = link_2
	ws.ServerId = tcp.ServerId

	go tcp.Run()
	go ws.Run()

	http.HandleFunc("/username", ws.Username)
	http.HandleFunc("/", ws.Home)

	log.Fatal(http.ListenAndServe(wsPort, nil))
}

package follower

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yoshifrancis/go-gameserver/internal/follower/tcpserver"
	"github.com/yoshifrancis/go-gameserver/internal/follower/wsserver"
)

type Follower struct {
	tcp *tcpserver.TCPServer
	ws  *wsserver.WSServer
}

func Follower_init(wsPort, leaderIp string, done chan bool) *Follower {

	tcp := tcpserver.ConnectToLeader(leaderIp, done)
	fmt.Println("Listening for WS connections at: ", wsPort)
	if tcp == nil {
		fmt.Println("Error connecting to leader")
		return &Follower{nil, nil}
	}
	fmt.Println("Connected to leader")

	ws := wsserver.NewWSServer(done)
	link_1 := make(chan []byte)
	link_2 := make(chan wsserver.LeaderRequest)
	tcp.WSfrom = link_1
	ws.TCPto = link_1
	tcp.WSto = link_2
	ws.TCPfrom = link_2
	ws.ServerId = tcp.ServerId

	go tcp.Run()
	go ws.Run()

	http.HandleFunc("/home", ws.Home)
	http.HandleFunc("/ping", ws.Ping)
	http.HandleFunc("/", ws.Index)

	defer log.Fatal(http.ListenAndServe(wsPort, nil))

	return &Follower{tcp, ws}
}

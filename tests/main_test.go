package tests_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/internal/follower"
	"github.com/yoshifrancis/go-gameserver/internal/leader"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

const (
	leaderTcpPort  = ":8000"
	followerWsPort = ":8001"
)

var upgrader = websocket.Upgrader{}

func TestLeaderInit(t *testing.T) {
	leader.Leader_init(leaderTcpPort, nil)

	conn, err := net.Dial("tcp", "localhost"+leaderTcpPort)
	if err != nil {
		t.Fatalf("Error connecting to the leader! Got the error: %s", err)
	}

	_, err = conn.Write([]byte(messages.Ping()))
	if err != nil {
		t.Fatalf("Error pinging the leader! Got the error: %s", err)
	}

	pongBuffer := make([]byte, 24)

	_, err = conn.Read(pongBuffer)

	if err != nil {
		t.Fatalf("Error receiving pong from the leader! Got the error: %s", err)
	}

	if string(pongBuffer) != messages.Pong() {
		t.Fatalf("Did not receive pong message! Instead received: %s", string(pongBuffer))
	}

}

// should be able to connect to leader
func TestFollowerInit(t *testing.T) {
	leader.Leader_init(leaderTcpPort, nil)
	follower.Follower_init(followerWsPort, leaderTcpPort, nil)
	var tests = struct {
		name  string
		url   string
		input string
		want  string
	}{
		name:  "Basic ping and pong to follower websocket",
		url:   "ws://localhost" + followerWsPort,
		input: messages.Ping(),
		want:  messages.Pong(),
	}
	for i := range 10 {
		t.Run(tests.name, func(t *testing.T) {
			ws, _, err := websocket.DefaultDialer.Dial(tests.url, nil)
			if err != nil {
				t.Fatalf("Unable to connect to follower websocket server! Got the error: %s", err)
			}
			defer ws.Close()

			if err := ws.WriteMessage(websocket.TextMessage, []byte(tests.input)); err != nil {
				t.Errorf("Error pinging to websocket server! Got the error: %s", err)
				return
			}
			_, buffer, err := ws.ReadMessage()

			if err != nil {
				t.Errorf("Error reading pong from websocket server! Got the error: %s", err)
				return
			}

			if string(buffer) != tests.want {
				t.Errorf("Incorrect pong message on iteration %d. Got %s / Expected %s", i, string(buffer), tests.want)
				return
			}
		})
	}
}

func TestFollowerLogin(t *testing.T) {
	leader.Leader_init(leaderTcpPort, nil)
	follower.Follower_init(followerWsPort, leaderTcpPort, nil)
	url := "ws://localhost" + followerWsPort
	type args struct {
		username string
		message  string
	}
	var tests = []struct {
		name  string
		input args
		want1 int
		want2 int
	}{
		{name: "First name and sending message", input: args{username: "Yoshi", message: "Is King"}, want1: 200, want2: 200},
		{name: "Second user and message", input: args{username: "Mario", message: "Is King"}, want1: 200, want2: 200},
		{name: "Testing Ping", input: args{username: messages.Ping(), message: messages.Ping()}, want1: 202, want2: 202},
		{name: "Third user and message", input: args{username: "Luigi", message: "Is King"}, want1: 200, want2: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatalf("Unable to connect to follower websocket server! Got the error: %s", err)
			}
			defer ws.Close()

			if err := ws.WriteMessage(websocket.TextMessage, []byte(tt.input.username)); err != nil {
				t.Errorf("Error sending username to websocket server! Got the error: %s", err)
				return
			}
			_, buffer, err := ws.ReadMessage()

			if err != nil {
				t.Errorf("Error reading status after sending username! Got the error: %s", err)
				return
			}

			status, _ := strconv.Atoi(string(buffer))
			if status != tt.want1 {
				t.Errorf("Incorrect status after sending username. Got %d / Expected %d", status, tt.want1)
				return
			}

			if err := ws.WriteMessage(websocket.TextMessage, []byte(tt.input.message)); err != nil {
				t.Errorf("Error sending message websocket server! Got the error: %s", err)
				return
			}
			_, buffer, err = ws.ReadMessage()

			if err != nil {
				t.Errorf("Error status after sending message! Got the error: %s", err)
				return
			}

			status, _ = strconv.Atoi(string(buffer))
			if status != tt.want2 {
				t.Errorf("Incorrect status after sending message. Got %d / Expected %d", status, tt.want2)
				return
			}
		})
	}

}

func TestLeaderShutdown(t *testing.T) {
	t.Run("Leader and Follower able to signal they shutdown", func(t *testing.T) {
		leaderDone := make(chan bool)
		followerDone := make(chan bool)
		leader := leader.Leader_init(leaderTcpPort, leaderDone)
		time.Sleep(time.Second) // ensure that leader gets set up (hopefully)
		follower.Follower_init(followerWsPort, leaderTcpPort, followerDone)

		go func() {
			time.Sleep(time.Second * 2)
			leader.Shutdown()
		}()

		timeToShutDown := time.NewTimer(time.Second * 4)

		select {
		case <-leaderDone:
			break
		case <-timeToShutDown.C:
			t.Fatal("Server did not shutdown within a reasonable amount of time!")
		}

		timeToShutDown = time.NewTimer(time.Second * 4)

		select {
		case <-followerDone:
			break
		case <-timeToShutDown.C:
			t.Fatal("Follower did not shutdown within a reasonable amount of time!")
		}
	})

}

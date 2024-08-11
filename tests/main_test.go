package tests_test

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yoshifrancis/go-gameserver/internal/follower"
	"github.com/yoshifrancis/go-gameserver/internal/leader"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

const (
	leaderTcpPort  = ":8000"
	leaderTcpAPort = ":8002"
	followerWsPort = ":8001"
)

func TestLeaderInit(t *testing.T) {
	t.Run("Initialize Server and Ping Through TCP", func(t *testing.T) {
		leader.Leader_init(leaderTcpPort, leaderTcpAPort, nil)

		conn, err := net.Dial("tcp", "localhost"+leaderTcpPort)
		if err != nil {
			t.Fatalf("Error connecting to the leader! Got the error: %s", err)
		}

		_, err = conn.Write([]byte(messages.Ping()))
		if err != nil {
			t.Fatalf("Error pinging the leader! Got the error: %s", err)
		}

		pongBuffer := make([]byte, 64)

		_, err = conn.Read(pongBuffer)
		fmt.Println("Got leader reply!", string(pongBuffer))

		if err != nil {
			t.Fatalf("Error receiving pong from the leader! Got the error: %s", err)
		}

		flag, pong := messages.Decode(pongBuffer)

		if flag != '!' && pong[0] != "PONG" {
			t.Fatalf("Did not receive pong message! Instead received: %s%s", string(flag), pong[0])
		}
	})

}

// should be able to connect to leader
func TestFollowerInit(t *testing.T) {
	fmt.Print("BEGINNING TESTS\n\n\n\n")

	leader.Leader_init(leaderTcpPort, leaderTcpAPort, nil)
	go follower.Follower_init(followerWsPort, leaderTcpPort, nil)

	fmt.Print("RUNNING TESTS\n\n\n\n")

	var tests = struct {
		name  string
		url   string
		input string
		want  string
	}{
		name:  "Basic ping and pong to follower websocket",
		url:   "ws://localhost" + followerWsPort + "/ping",
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
	leader.Leader_init(leaderTcpPort, "0", nil)
	go follower.Follower_init(followerWsPort, leaderTcpPort, nil)
	url := "ws://localhost" + followerWsPort + "/home"
	type args struct {
		username string
		message  string
	}
	var tests = []struct {
		name  string
		input args
		want1 string
	}{
		{name: "First name and sending message", input: args{username: "Yoshi", message: "Is King"}, want1: "<div id=\"app\" hx-swap=\"outerHTML\">"},
		{name: "Second user and message", input: args{username: "Mario", message: "Is King"}, want1: "<div id=\"app\" hx-swap=\"outerHTML\">"},
		{name: "Testing Ping", input: args{username: messages.Ping(), message: messages.Ping()}, want1: "<div id=\"app\" hx-swap=\"outerHTML\">"},
		{name: "Third user and message", input: args{username: "Luigi", message: "Is King"}, want1: "<div id=\"app\" hx-swap=\"outerHTML\">"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Fatalf("Unable to connect to follower websocket server! Got the error: %s", err)
			}
			defer ws.Close()

			_, _, err = ws.ReadMessage() // username prompt buffer

			if err != nil {
				t.Errorf("Error reading username prompt! Got the error: %s", err)
			}

			type Username struct {
				Username string `json:"username"`
			}

			jsonMsg, _ := json.Marshal(Username{Username: tt.input.username})
			if err := ws.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
				t.Errorf("Error sending username to websocket server! Got the error: %s", err)
				return
			}
			_, buffer, err := ws.ReadMessage()

			if err != nil {
				t.Errorf("Error reading status after sending username! Got the error: %s", err)
				return
			}

			if string(buffer[:34]) != tt.want1 {
				t.Errorf("Incorrect status after sending username. Got %s / Expected %s", string(buffer[:34]), tt.want1)
				return
			}
		})
	}

}

func TestLeaderShutdown(t *testing.T) {
	t.Run("Leader and Follower able to signal they shutdown", func(t *testing.T) {
		leaderDone := make(chan bool)
		followerDone := make(chan bool)
		leader := leader.Leader_init(leaderTcpPort, "0", leaderDone)
		go follower.Follower_init(followerWsPort, leaderTcpPort, followerDone)

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
		fmt.Println("Leader finished!")
		timeToShutDown = time.NewTimer(time.Second * 4)
		select {
		case <-followerDone:
			break
		case <-timeToShutDown.C:
			t.Fatal("Follower did not shutdown within a reasonable amount of time!")
		}
	})

}

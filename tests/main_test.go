package tests_test

import (
	"net"
	"testing"

	"github.com/yoshifrancis/go-gameserver/internal/follower"
	"github.com/yoshifrancis/go-gameserver/internal/leader"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

const (
	leaderTcpPort  = ":8000"
	followerWsPort = ":8001"
)

func TestLeaderInit(t *testing.T) {
	leader.Leader_init(leaderTcpPort)

	conn, err := net.Dial("tcp", leaderTcpPort)
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
	follower.Follower_init(followerWsPort, leaderTcpPort)
}

func TestLeaderShutdown(t *testing.T) {

}

func TestFollowerShutdown(t *testing.T) {

}

func TestSendMessage(t *testing.T) {

}

func TestJoinMultipleUsers(t *testing.T) {

}

func TestLeaveUser(t *testing.T) {

}

func TestJoinLobby(t *testing.T) {

}

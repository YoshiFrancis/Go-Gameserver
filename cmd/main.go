package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/yoshifrancis/go-gameserver/internal/follower"
	"github.com/yoshifrancis/go-gameserver/internal/leader"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <TYPE OF SERVER (L, F)> <PORT>")
	}

	typing := os.Args[1]
	port := os.Args[2]
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Given an invalid port num")
	}

	if typing == "L" {
		leader.Leader_init(":"+port, nil)
	} else if typing == "F" {
		fmt.Println("What is the address of the leader you would like to follow?")
		leaderip := input()
		follower.Follower_init(":"+port, leaderip, nil)
	}

	for {
		user_input := input()
		if user_input == "exit" {
			break
		}
	}
}

func input() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadBytes('\n')
	text = []byte(strings.Replace(string(text), "\n", "", -1))
	text_string := string(text)
	return text_string
}

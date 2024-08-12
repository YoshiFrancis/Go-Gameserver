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
	if len(os.Args) != 3 && len(os.Args) != 4 {
		log.Fatal("Usage: go run main.go <TYPE OF SERVER (L, F)> <PORT> <PORT if L")
	}

	typing := os.Args[1]
	port := os.Args[2]
	var port2 string
	if len(os.Args) == 4 {
		port2 = os.Args[3]
	}
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("Given an invalid port num")
	}
	_, err = strconv.Atoi(port2)
	if err != nil {
		fmt.Println("Given an invalid port num")
	}

	if typing == "L" && len(os.Args) == 4 {
		fmt.Println("in here")
		leader.Leader_init(":"+port, ":"+port2, nil)
	} else if typing == "F" && len(os.Args) == 3 {
		fmt.Println("What is the address of the leader you would like to follow?")
		leaderip := input()
		go follower.Follower_init(":"+port, leaderip, nil)
	} else {
		return
	}

	for {
		fmt.Println("HELLOs")
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

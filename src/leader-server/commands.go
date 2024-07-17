package leaderserver

import (
	"fmt"
	"strings"
)

func handleArgs(flag byte, args []string) (res string) {
	args[0] = strings.ToLower(args[0])
	res = ""
	if flag == '-' {
		switch args[0] {
		case "create":
			break
		case "shutdown":
			break
		case "disc":
			break
		case "accept":
			break
		default:
			fmt.Println("Given an invalid server command")
			return
		}

	}

	return ""
}

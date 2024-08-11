package messages

import "fmt"

func ApplicationBroadcast(arg, lobbyTitle string, usernames []string) string {
	receiversStr := listUsernames(usernames)
	message := fmt.Sprintf("_4\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(receiversStr), receiversStr)
	return message
}

func ForApplicationRequest(arg, lobbyTitle, sender string) string {
	message := fmt.Sprintf("_4\r\n7\r\nREQUEST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(sender), sender)
	return message
}

func ApplicationShutdown() string {
	message := "_1\r\n8\r\nSHUTDOWN\r\n\r\n"
	return message
}

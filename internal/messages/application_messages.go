package messages

import (
	"fmt"
	"strings"
)

type ApplicationRequest struct {
	Flag       byte
	Command    string
	Arg        string
	LobbyTitle string
	Sender     string
	Receivers  []string
}

func AReqDecode(req []byte) ApplicationRequest {
	flag, args := Decode(req)
	if flag == 'x' {
		return ApplicationRequest{Flag: 'x'}
	}
	fmt.Println(args)
	args[0] = strings.ToLower(args[0])
	if len(args) == 1 && args[0] == "shutdown" {
		return ApplicationRequest{
			Flag:    flag,
			Command: args[0],
		}
	} else if len(args) == 5 && (args[0] == "send" || args[0] == "request" || args[0] == "start") {
		usernames := unlistUsernames(args[4])
		return ApplicationRequest{
			Flag:       flag,
			Command:    strings.ToLower(args[0]),
			Arg:        args[1],
			LobbyTitle: args[2],
			Sender:     args[3],
			Receivers:  usernames,
		}
	} else if len(args) == 4 {
		return ApplicationRequest{
			Flag:       flag,
			Command:    strings.ToLower(args[0]),
			Arg:        args[1],
			LobbyTitle: args[2],
			Sender:     args[3],
			Receivers:  []string{},
		}
	} else {
		return ApplicationRequest{
			Flag: 'x',
		}
	}
}

func ApplicationBroadcast(arg, lobbyTitle, sender string) string { // directed to whole lobby
	message := fmt.Sprintf("_4\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(sender), sender)
	return message
}

func ApplicationStart(arg, lobbyTitle, sender string, usernames []string) string {
	receiversStr := listUsernames(usernames)
	message := fmt.Sprintf("_5\r\n5\r\nSTART\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(sender), sender, len(receiversStr), receiversStr)
	return message
}

func ApplicationSend(arg, lobbyTitle, sender string, usernames []string) string { // directed to particular uses
	receiversStr := listUsernames(usernames)
	message := fmt.Sprintf("_5\r\n4\r\nSEND\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(sender), sender, len(receiversStr), receiversStr)
	return message
}

func ApplicationRequestTo(arg, lobbyTitle, sender string, usernames []string) string { // send to the application server from leader
	receiversStr := listUsernames(usernames) // we send usernamees in case new users joined lobby, we let application handle such
	message := fmt.Sprintf("_5\r\n7\r\nREQUEST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(lobbyTitle), lobbyTitle, len(sender), sender, len(receiversStr), receiversStr)
	return message
}

func ApplicationShutdown() string {
	message := "_1\r\n8\r\nSHUTDOWN\r\n\r\n"
	return message
}

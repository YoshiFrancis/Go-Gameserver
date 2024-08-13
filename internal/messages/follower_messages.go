package messages

import (
	"fmt"
)

func DisconnectUser(username string) string {
	message := fmt.Sprintf("-2\r\n4\r\nDISC\r\n%d\r\n%s\r\n\r\n", len(username), username)
	return message
}

func RegisterUser(username string) string {
	message := fmt.Sprintf("-2\r\n8\r\nREGISTER\r\n%d\r\n%s\r\n\r\n", len(username), username)
	return message
}

func FollowerRoomBroadcast(broadcast, username string) string {
	message := fmt.Sprintf("-3\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(broadcast), broadcast, len(username), username)
	return message
}

func RoomJoinUser(lobbyTitle, username string) string {
	message := fmt.Sprintf("+3\r\n4\r\nJOIN\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(lobbyTitle), lobbyTitle, len(username), username)
	return message
}

func CreateLobby(lobbyTitle, username string) string {
	message := fmt.Sprintf("+3\r\n5\r\nLOBBY\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(lobbyTitle), lobbyTitle, len(username), username)
	return message
}
func FollowerAppStart(appName, username string) string {
	message := fmt.Sprintf("+3\r\n9\r\nAPP-START\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(appName), appName, len(username), username)
	return message
}

func FollowerAppRequest(arg, username string) string {
	message := fmt.Sprintf("+3\r\n3\r\nAPP\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(arg), arg, len(username), username)
	return message
}

func UserMessageUser(origin_username string, target_username string, private_message string) string {
	message := fmt.Sprintf("*4\r\n2\r\nPM\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(origin_username), origin_username, len(target_username), target_username, len(private_message), private_message)
	return message
}

package messages

import (
	"fmt"
	"strconv"
)

func DisconnectUser(username string) string {
	message := fmt.Sprintf("-2\r\n4\r\nDISC\r\n%d\r\n%s\r\n\r\n", len(username), username)
	return message
}

func RegisterUser(username string, serverId int) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-3\r\n8\r\nREGISTER\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(username), username)
	return message
}

func FollowerRoomBroadcast(username string, roomId int, broadcast string) string {
	roomIdStr := strconv.Itoa(roomId)
	roomIdLength := len(roomIdStr)
	message := fmt.Sprintf("+4\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", roomIdLength, roomIdStr, len(username), username, len(broadcast), broadcast)
	return message
}

func RoomJoinUser(username string, roomId int) string {
	roomIdStr := strconv.Itoa(roomId)
	roomIdLength := len(roomIdStr)
	message := fmt.Sprintf("+3\r\n4\r\nJOIN\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", roomIdLength, roomIdStr, len(username), username)
	return message
}

func CreateLobby(lobbyTitle, username string) string {
	message := fmt.Sprintf("+3\r\n5\r\nLOBBY\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(username), username, len(lobbyTitle), lobbyTitle)
	return message
}

func UserMessageUser(origin_username string, target_username string, private_message string) string {
	message := fmt.Sprintf("*4\r\n2\r\nPM\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(origin_username), origin_username, len(target_username), target_username, len(private_message), private_message)
	return message
}

package messages

import (
	"fmt"
	"strconv"
)

// - means request going to follower, so send HTMX
// + means request going to leader, so send in REDIS form

func LeaderAcceptServer(serverId int, url string) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n6\r\nACCEPT\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(url), url)
	return message
}

func LeaderCreation(serverId int, url string) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n8\r\nCREATION\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(url), url)
	return message
}

func LeaderShutdown(serverId int) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n8\r\nSHUTDOWN\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr)
	return message
}

func ServerMergeData(hubId int, lobbyIds []int, userData [][]string, idGenStartingInt int) string {
	hubIdStr := strconv.Itoa(hubId)
	hubIdLength := len(hubIdStr)
	var lobbyIdsString string
	for lobbyId := range lobbyIds {
		lobbyIdStr := strconv.Itoa(lobbyId)
		lobbyIdStrLength := len(lobbyIdStr)
		new_lobby := fmt.Sprintf("%d\r\n%s\r\n", lobbyIdStrLength, lobbyIdStr)
		lobbyIdsString += new_lobby
	}
	lobbyIdsStringLength := len(lobbyIdsString)

	var userDataStr string
	for _, user := range userData {
		username := user[0]
		serverId := user[1]
		roomId := user[2]
		new_user := fmt.Sprintf("%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n", len(username), username, len(serverId), serverId, len(roomId), roomId)
		userDataStr += new_user
	}
	userDataStrLength := len(userDataStr)

	idGenStr := strconv.Itoa(idGenStartingInt)
	idGenStrLength := len(idGenStr)

	total_args := 1 + 1 + len(lobbyIds) + len(userData) + 1 // the command, hub Id, total lobbies, total users, id gen starting int
	message := fmt.Sprintf("-%d\r\n5\r\nMERGE\r\n%d\r\n%s\r\n%d\r\n%s%d\r\n%s\r\n%d\r\n%s\r\n\r\n", total_args, hubIdLength, hubIdStr, lobbyIdsStringLength, lobbyIdsString, userDataStrLength, userDataStr, idGenStrLength, idGenStr)
	return message
}

func LeaderRoomBroadcast(sender, broadcast string, usernames []string) string {
	usernamesStr := listUsernames(usernames)
	message := fmt.Sprintf("+4\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(sender), sender, len(broadcast), broadcast, len(usernamesStr), usernamesStr)
	return message
}

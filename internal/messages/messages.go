package messages

import (
	"bytes"
	"fmt"
	"strconv"
)

/*
MESSAGE STILL NEED:
merge user
merge server
tell user to join a specific url
*/

func Decode(req []byte) (flag byte, args []string) {
	r := bytes.NewReader(req)
	flag, _ = r.ReadByte() // reading the flag. lowkey dont know what to do with it right now
	argcByte, _ := r.ReadByte()
	argc := argcByte - '0'
	args = make([]string, int(argc))

	if flag != '-' && flag != '+' && flag != '/' && flag != '!' {
		flag = 'x'
		return
	}

	readCLRF(r)

	for idx := range argc {
		sizeInt, ok := readSize(r)
		if !ok {
			fmt.Println("Read an invalid size")
		}
		var currArg []byte
		for range sizeInt {
			b, _ := r.ReadByte()
			currArg = append(currArg, b)
		}
		args[idx] = string(currArg)
		if !readCLRF(r) {
			flag = 'x' // signal invalid message
			return
		}
	}

	return flag, args
}

func readCLRF(r *bytes.Reader) bool {
	if b, err := r.ReadByte(); b != '\r' || err != nil {
		fmt.Println("Error reading CLRF while decoding")
		return false
	} // '\r

	if b, err := r.ReadByte(); b != '\n' || err != nil {
		fmt.Println("Error reading CLRF while decoding")
		return false
	} // '\n
	return true
}

func readSize(r *bytes.Reader) (int, bool) {
	var size []byte

	for {
		b, _ := r.ReadByte()
		if b == '\r' {
			break
		}
		size = append(size, b)
	}
	if b, err := r.ReadByte(); b != '\n' || err != nil {
		return -1, false
	}

	sizeInt, err := strconv.Atoi(string(size))
	if err != nil {
		return -1, false
	}
	return sizeInt, true
}

// - -> for server
// + -> for Hub
// / -> for rooms
// * -> for individual clients
// _ -> for groups (if implemented in the future)
// flag \r\n n arguments \r\n arg[0] \r\n arg[1] ... \r\n arg[n-1] \r\n\r\n

func Ping() string {
	message := "!1\r\n4\r\nPING\r\n\r\n"
	return message
}

func Pong() string {
	message := "!1\r\n4\r\nPONG\r\n\r\n"
	return message
}

func ServerAcceptServer(serverId int, url string) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n6\r\nACCEPT\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(url), url)
	return message
}

func ServerCreation(serverId int, url string) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n8\r\nCREATION\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(url), url)
	return message
}

func ServerShutdown(serverId int) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n8\r\nSHUTDOWN\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr)
	return message
}

func ServerTellServerId(serverId int) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-2\r\n8\r\nSERVERID\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr)
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

func ServerDisconnectUser(username string) string {
	message := fmt.Sprintf("-2\r\n4\r\nDISC\r\n%d\r\n%s\r\n\r\n", len(username), username)
	return message
}

func ServerRegisterUser(username string, serverId int) string {
	serverIdStr := strconv.Itoa(serverId)
	serverIdLength := len(serverIdStr)
	message := fmt.Sprintf("-3\r\n8\r\nREGISTER\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", serverIdLength, serverIdStr, len(username), username)
	return message
}

func RoomBroadcast(username string, roomId int, broadcast string) string {
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

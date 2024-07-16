package messages

import (
	"bytes"
	"fmt"
	"strconv"
)

func Decode(req []byte) (flag byte, args []string) {
	r := bytes.NewReader(req)
	flag, _ = r.ReadByte() // reading the flag. lowkey dont know what to do with it right now
	argc, _ := r.ReadByte()
	args = make([]string, int(argc))

	if flag != '-' && flag != '+' && flag != '*' {
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
	for b, _ := r.ReadByte(); b != '\r'; {
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
// + -> for rooms
// * -> for individual clients
// _ -> for groups (if implemented in the future)
// flag \r\n n arguments \r\n arg[0] \r\n arg[1] ... \r\n arg[n-1] \r\n\r\n

func ServerShutdown() string {
	message := "-1\r\n8\r\nSHUTDOWN\r\n\r\n"
	return message
}

func ServerDisconnectUser(username string, roomId int) string {
	roomIdStr := strconv.Itoa(roomId)
	roomIdLength := len(roomIdStr)
	message := fmt.Sprintf("-3\r\n4\r\nDISC\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", roomIdLength, roomIdStr, len(username), username)
	return message
}

func RoomJoinUser(username string, roomId int) string {
	roomIdStr := strconv.Itoa(roomId)
	roomIdLength := len(roomIdStr)
	message := fmt.Sprintf("+3\r\n4\r\nJOIN\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", roomIdLength, roomIdStr, len(username), username)
	return message
}

func RoomBroadcast(username string, roomId int, broadcast string) string {
	roomIdStr := strconv.Itoa(roomId)
	roomIdLength := len(roomIdStr)
	message := fmt.Sprintf("+4\r\n9\r\nBROADCAST\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", roomIdLength, roomIdStr, len(username), username, len(broadcast), broadcast)
	return message
}

func UserMessageUser(origin_username string, target_username string, private_message string) string {
	message := fmt.Sprintf("*4\r\n2\r\nPM\r\n%d\r\n%s\r\n%d\r\n%s\r\n%d\r\n%s\r\n\r\n", len(origin_username), origin_username, len(target_username), target_username, len(private_message), private_message)
	return message
}

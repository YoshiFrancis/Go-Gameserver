package messages

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

/*
MESSAGE STILL NEED:
merge user
merge server
tell user to join a specific url
*/

type FollowerRequest struct {
	Flag    byte
	Command string
	Sender  string
	Arg     string
}

type LeaderRequest struct {
	Flag      byte
	Command   string
	Arg       string
	Receivers []string
}

// follower request format
// command name, argument, username of sender,
// some follower commands do not have a sender

func FReqDecode(req []byte) FollowerRequest {
	flag, args := Decode(req)
	fReq := FollowerRequest{Flag: flag}
	fReq.Command = strings.ToLower(args[0])
	fReq.Arg = args[1]
	if len(args) > 2 {
		fReq.Sender = args[2]
	}
	return fReq
}

// leader request format
// command name, argument, list of usernames sepearated by \n

func LReqDecode(req []byte) LeaderRequest {
	flag, args := Decode(req)
	fmt.Println("Decoded: ", args)
	usernames := unlistUsernames(args[2])
	return LeaderRequest{
		Flag:      flag,
		Command:   strings.ToLower(args[0]),
		Arg:       args[1],
		Receivers: usernames,
	}
}

func Decode(req []byte) (byte, []string) {

	r := bytes.NewReader(req)
	flag, _ := r.ReadByte() // reading the flag. lowkey dont know what to do with it right now
	argcByte, _ := r.ReadByte()
	argc := argcByte - '0'
	args := make([]string, int(argc))

	if flag != '-' && flag != '+' && flag != '/' && flag != '!' && flag != '_' {
		return 'x', []string{}
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
			return 'x', []string{}
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

func listUsernames(usernames []string) string {
	list := ""
	for _, username := range usernames {
		list += username + "\n"
	}

	return list
}

func unlistUsernames(usernameStr string) []string {
	usernames := make([]string, 0)
	currName := make([]byte, 0)
	r := bytes.NewReader([]byte(usernameStr))
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if b == '\n' {
			usernames = append(usernames, string(currName))
			currName = make([]byte, 0)
		} else {
			currName = append(currName, b)
		}
	}

	return usernames
}

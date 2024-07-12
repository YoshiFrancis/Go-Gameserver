package messages

import "fmt"

type Message struct {
	body        string
	body_length int
	flag        rune
}

func decode(msg []byte) string {
	fmt.Println(string(msg))
	return ""
}

func encode(msg string) string {
	fmt.Println(msg)
	return ""
}

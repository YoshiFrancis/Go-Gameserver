/*

Assumption
1. Everything will be linear, so no use of channels or go routines in order to reduce complexity
2. Only one leader and once the leader shuts down, application does as well

Preliminary
1. Settup html file
2. Decide how you want input and configure the html file to send it in the right format (may have to use jQuery and events to do for some actions, like a mouse click)


Run-time
1. Connect to leader (or set up a listener connection and listen for leader connections)
2. Wait for confirmation to start spplication, with a list of receivers
3. Separate users into a container of some sorts (like a room), and then render htmx template and send back to leader in
		messages.FollowerBroadcast format (users can be separated via their lobby title)
4. Extra usernames may trinkle in, handle them (deny or accept into application)
5. Handle any input that comes in from the server
5. When server signals shutdown, close connection to that leader server
*/

package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"os"

	"github.com/yoshifrancis/go-gameserver/internal/containers"
	"github.com/yoshifrancis/go-gameserver/internal/messages"
)

var appTemplate *template.Template

type ChatRoom struct {
	usernames  []string
	messages   []string
	lobbyTitle string
}

func init() {
	var err error
	appTemplate, err = template.ParseFiles("app.html")
	if err != nil {
		fmt.Println("Error parsing html file!")
		os.Exit(1)
	}
}

func main() {
	appName := "chatroom"
	addr := ":8002"

	appRooms := make(map[string]ChatRoom) // lobby title to list of usernames

	conn := connectToLeader(addr, appName)

	defer conn.Close()

	for {
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error reading from leader. Closing connection!")
			panic(err) // panicing since this is only leader
		}

		// decode buffer using messages package
		aReq := messages.AReqDecode(buff[:n])

		var res []byte // variable to hold responses

		// only can receive three commands from leader: start, request, shutdown

		// if command is to start, create chat room
		if aReq.Command == "start" {
			fmt.Println("received singal to start chatroom")
			chatRoom := ChatRoom{
				usernames:  aReq.Receivers,
				messages:   []string{},
				lobbyTitle: aReq.LobbyTitle,
			}
			appRooms[aReq.LobbyTitle] = chatRoom

			// send out the application template
			tmpl := containers.RenderTemplate(appTemplate, struct {
				Participants []string
				Messages     []string
			}{
				Participants: chatRoom.usernames,
				Messages:     []string{},
			})

			str_res := messages.ApplicationBroadcast(string(tmpl), chatRoom.lobbyTitle, "Server")
			res = []byte(str_res)

		} else if aReq.Command == "request" { // else if command is a request, let chat room handle request
			c, ok := appRooms[aReq.LobbyTitle]
			if !ok {
				log.Println("Received invalid lobby title!")
				continue
			}
			res = c.addMessage(aReq.Arg, aReq.Sender)
		} else if aReq.Command == "shutdown" { // else if command is a shutdown, close connection to leader and shutdown application
			return
		} else { // if none of the above command, do nothing and continue on to next read
			log.Println("Received unknown command: ", aReq.Command) // log it for debugging purposes
			log.Println("Total received buffer: ", string(buff[:n]))
			continue
		}

		// send back handled response
		log.Println("Sending back response: ", string(res))
		conn.Write(res)
	}
}

func connectToLeader(addr, appName string) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to leader!")
		panic(err) // if failed to connect to leader, crash program. For more robust programs where you may connect to several leaders, find better ways
	}

	_, err = conn.Write([]byte(appName)) // send application name

	if err != nil {
		fmt.Println("Error sending application name to leader!")
		panic(err)
	}

	return conn
}

// adding message to the message queue currently held by the chat room
// templating the message with the sender into a <li> format for the html
// placing all the <li> tags into a <ul> with an id="app-messages" because with htmx websockets,
// you replace element using out of bounds swapping, which swaps elements via the id
func (c *ChatRoom) addMessage(msg, sender string) []byte {
	c.messages = append(c.messages, templateMessage(sender, msg))
	broadcastMsg := "<ul id=\"app-messages\">"
	for _, message := range c.messages {
		broadcastMsg += message
	}
	broadcastMsg += "</ul>"

	return []byte(broadcastMsg)
}

func templateMessage(sender, msg string) string {
	return "<li>" + sender + ": " + msg + "</li>"
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"

	"kelber.com/connect4/msg"
)

const (
	connHost = "localhost"
	connPort = "8080"
	connType = "tcp"
)

func main() {
	fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)

	conn, err := net.Dial(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

	go listen(&conn)
	for {
		username, _ := reader.ReadString('\n')
		m := msg.CreateNewMessage(msg.Request, msg.NewPlayerReq, msg.DefaultContentDelimiter, [][]byte{[]byte(username)})
		buffer := bytes.Buffer{}
		msg.Serialize(m, &buffer)
		fmt.Println("Creating message: ", buffer.Bytes())
		conn.Write(buffer.Bytes())

	}
}

func listen(conn *net.Conn) {
	for {
		fmt.Println("Listening for response")
		b, err := bufio.NewReader(*conn).ReadBytes(msg.EndOfMessage)
		if err != nil {
			fmt.Println("Server is down.")
			(*conn).Close()
			return
		}
		buffer := bytes.NewBuffer(b)
		message := msg.Deserialize(bytes.Buffer(*buffer))
		fmt.Println("received message: ", message)
		parseResponse(conn, message)
	}
}

func parseResponse(conn *net.Conn, message msg.Message) {
	if message.Type != msg.Response {
		fmt.Printf("Failed to parse response with message type %d\n", message.Type)
		return
	}
	switch message.ID {
	case msg.NewPlayerResp:

	case msg.StartGameResp:

	case msg.StartTurnResp:

	case msg.PlacePieceResp:

	case msg.UpdateLobbyResp:
		UpdateState(conn, message)

	default:
		fmt.Printf("Received message: %v\n", message)

	}
}

func UpdateState(conn *net.Conn, message msg.Message) {
	fmt.Println("Update state...")
	builder := strings.Builder{}
	for _, username := range message.Content {
		builder.Write(username)
		builder.WriteString(" ")
	}
	fmt.Println(builder.String())
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"

	"kelber.com/connect4/game"
	"kelber.com/connect4/msg"
)

const (
	connHost = "127.0.0.1"
	connPort = "8080"
	connType = "tcp"
)

var lobby game.Lobby

func main() {
	lobby = game.CreateLobby()
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
		// m := msg.CreateNewMessage(msg.Response, msg.StartTurnResp, []byte{1})
		// fmt.Println(m)
		// fmt.Print("Serializing into: ")
		// buffer := bytes.Buffer{}
		// msg.Serialize(m,&buffer)
		// fmt.Printf("%v\n", buffer)

		// fmt.Print("Deserializing back to: ")
		// mm := msg.Deserialize(buffer)
		// fmt.Printf("%v\n", mm)
		go handleConnection(&c)
	}
}

func handleConnection(conn *net.Conn) {
	b, err := bufio.NewReader(*conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("Client left.")
		(*conn).Close()
		return
	}
	buffer := bytes.NewBuffer(b)
	message := msg.Deserialize(bytes.Buffer(*buffer))

	go func() {
		parseMessage(conn, message)
		handleConnection(conn)
	}()
}

func parseMessage(conn *net.Conn, message msg.Message) bool {
	switch message.Type {
	case msg.Request:
		return handleRequest(conn, message)
	case msg.Response:
		return handleResponse(conn, message)
	default:
		fmt.Printf("Unknown message type: %d\n", message.Type)
		return false
	}

}

func handleRequest(conn *net.Conn, message msg.Message) bool {
	switch message.ID {
	case msg.NewPlayerReq:
		CreateNewPlayer(conn, message)
	case msg.StartGameReq:

	case msg.StartTurnReq:
		player := int(message.Content[0])
		StartTurn(player)
	case msg.PlacePieceReq:

	case msg.UpdateStateReq:

	case msg.AnimationDoneReq:

	case msg.GameOverReq:

	}
	return true
}

func CreateNewPlayer(conn *net.Conn, message msg.Message) {
	username, err := bytes.NewBuffer(message.Content).ReadString('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Printf("Encountered error when reading player username %s\n", err)
		}
	}
	p := game.CreatePlayer(username, conn)
	lobby.AddPlayer(p)
	fmt.Printf("New lobby: %s\n", lobby)
}

func StartTurn(player int) {
	fmt.Printf("Player %d is starting their turn...\n", player)
}

func handleResponse(conn *net.Conn, message msg.Message) bool {
	switch message.ID {
	case msg.NewPlayerResp:

	case msg.StartGameResp:

	case msg.StartTurnResp:

	case msg.PlacePieceResp:

	case msg.UpdateStateResp:

	case msg.AnimationDoneResp:

	case msg.GameOverResp:

	}
	return true
}

package main

import (
	"fmt"
	"net"
	"os"

	"kelber.com/connect4/game"
)

const (
	connHost = "127.0.0.1"
	connPort = "8080"
	connType = "tcp"
)

//Maps IP address/Port to session
var activeSessions map[string]*game.Session

func main() {
	activeSessions = make(map[string]*game.Session)

	lobby := game.CreateLobby()
	requestHandler := game.CreateRequestHandler(&lobby, &activeSessions)

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

		// Create a new session
		// start a coroutine with the handle connection
		clientSession := game.CreateSession(&c, &requestHandler)
		activeSessions[clientSession.GetAddress()] = clientSession
		go clientSession.Listen()
	}
}

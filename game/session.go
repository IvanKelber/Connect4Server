package game

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"kelber.com/connect4/msg"
)

type Session struct {
	conn           *net.Conn
	requestHandler *RequestHandler
	Player         *Player
}

func (session *Session) Listen() {
	for {
		fmt.Println("Listening for new message")
		b, err := bufio.NewReader(*session.conn).ReadBytes('\n')
		if err != nil {
			session.OnExit()
			return
		}
		buffer := bytes.NewBuffer(b)
		message := msg.Deserialize(bytes.Buffer(*buffer))
		go func() {
			session.requestHandler.ParseMessage(session, message)
		}()
	}

}

func CreateSession(conn *net.Conn, requestHandler *RequestHandler) *Session {
	return &Session{conn, requestHandler, nil}
}

func (session *Session) SetPlayer(username string) {
	if session.Player == nil {
		session.Player = CreatePlayer(username)
	} else {
		session.Player.Username = username
	}
}

func (session *Session) OnExit() {

	session.requestHandler.DisconnectClient(session)

	fmt.Println("Client left.")
	(*session.conn).Close()
	return
}

func (session *Session) Write(bytes []byte) {
	fmt.Println("Writing bytes: ", bytes)
	(*session.conn).Write(bytes)
}

func (session *Session) GetAddress() string {
	return (*session.conn).RemoteAddr().String()
}

func (session *Session) GetUsername() string {
	if session.Player == nil {
		return ""
	}
	return session.Player.Username
}

func (session *Session) TestClientHandler() {
	var testByteStream = []byte{1, 28, 4, 28, 29, 28, 105, 118}
	var testByteStream2 = []byte{97, 110, 110, 110, 29, 28, 31}
	session.Write(testByteStream)
	session.Write(testByteStream2)
}

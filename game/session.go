package game

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"

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

func (session *Session) Write(buffer bytes.Buffer) {
	fmt.Println("Writing bytes: ", buffer.Bytes())
	(*session.conn).Write(buffer.Bytes())
}

func (session *Session) SendMessage(message msg.Message) {
	b := bytes.Buffer{}
	msg.Serialize(message, &b)
	session.Write(b)
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

//Sends a challenge request from the other user
func (session *Session) SendChallengeProposal(opponent string) {
	challenger := [][]byte{[]byte(opponent)}
	message := msg.CreateNewMessage(msg.Response, msg.ChallengeProposalResp, msg.DefaultContentDelimiter, challenger)
	session.SendMessage(message)
}

func (session *Session) WaitForChallengeResponse(opponent string) {
	content := [][]byte{[]byte(opponent)}

	message := msg.CreateNewMessage(msg.Response,
		msg.WaitForChallengeResp,
		msg.DefaultContentDelimiter,
		content)
	session.SendMessage(message)
}

func (session *Session) ChallengeRejected() {
	content := make([][]byte, 0)

	message := msg.CreateNewMessage(msg.Response,
		msg.ChallengeRejectedResp,
		msg.DefaultContentDelimiter,
		content)
	session.SendMessage(message)
}

func (session *Session) StartGame(g *Game) {
	content := make([][]byte, 0)

	var myTurn byte

	if g.IsMyTurn(session) {
		myTurn = msg.TrueByte
	} else {
		myTurn = msg.FalseByte
	}
	fmt.Println("GameId ", []byte(g.Id))
	content = append(content, []byte(g.Id))
	content = append(content, []byte{myTurn})

	message := msg.CreateNewMessage(msg.Response,
		msg.StartGameResp,
		msg.DefaultContentDelimiter,
		content)
	session.SendMessage(message)
}

func (session *Session) NotifyBoardUpdate(game *Game, column, winner int) {
	// Column, gameWon?, am I the winner?
	content := make([][]byte, 0)

	content = append(content, []byte(strconv.Itoa(column)))
	if winner == -1 {
		content = append(content, []byte{msg.FalseByte})
	} else {
		content = append(content, []byte{msg.TrueByte})
	}
	playerId, err := game.GetPlayerId(session)
	if err != nil {
		fmt.Println("Uh, player doesn't exist wtf")
		return
	}
	if winner == playerId {
		content = append(content, []byte{msg.TrueByte})
	} else {
		content = append(content, []byte{msg.FalseByte})
	}
	message := msg.CreateNewMessage(msg.Response,
		msg.PlacePieceResp,
		msg.DefaultContentDelimiter,
		content)
	session.SendMessage(message)

}

func (session *Session) TestClientHandler() {
	// var testByteStream = []byte{1, 28, 4, 28, 29, 28, 105, 118}
	// var testByteStream2 = []byte{97, 110, 110, 110, 29, 28, 31}
	// session.Write(testByteStream)
	// session.Write(testByteStream2)
}

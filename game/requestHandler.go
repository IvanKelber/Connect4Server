package game

import (
	"bytes"
	"fmt"
	"io"

	"kelber.com/connect4/msg"
)

type RequestHandler struct {
	lobby          *Lobby
	activeSessions *map[string]*Session
}

func CreateRequestHandler(lobby *Lobby, activeSessions *map[string]*Session) RequestHandler {
	return RequestHandler{lobby, activeSessions}
}

func (rh *RequestHandler) DisconnectClient(session *Session) {
	rh.lobby.Remove(session)
	delete(*rh.activeSessions, session.GetAddress())
	rh.lobby.UpdatePlayers()
}

func (rh *RequestHandler) ParseMessage(session *Session, message msg.Message) bool {
	switch message.Type {
	case msg.Request:
		return rh.handleRequest(session, message)
	case msg.Response:
		return rh.handleResponse(session, message)
	default:
		fmt.Printf("Unknown message type: %d\n", message.Type)
		return false
	}

}

func (rh *RequestHandler) handleRequest(session *Session, message msg.Message) bool {
	switch message.ID {
	case msg.NewPlayerReq:
		rh.CreateNewPlayer(session, message)
	case msg.StartGameReq:

	case msg.StartTurnReq:
		player := int(message.Content[0][0])
		rh.StartTurn(player)
	case msg.PlacePieceReq:

	case msg.UpdateStateReq:

	case msg.AnimationDoneReq:

	case msg.GameOverReq:

	}
	return true
}

func (rh *RequestHandler) CreateNewPlayer(session *Session, message msg.Message) {
	username, err := bytes.NewBuffer(message.Content[0]).ReadString('\n')
	if err != nil {
		if err != io.EOF {
			fmt.Printf("Encountered error when reading player username %s\n", err)
		}
	}

	if rh.lobby.AddPlayerSession(username, session) {
		session.SetPlayer(username)

		rh.lobby.UpdatePlayers()
	}
}

func (rh *RequestHandler) StartTurn(player int) {
	fmt.Printf("Player %d is starting their turn...\n", player)
}

func (rh *RequestHandler) handleResponse(session *Session, message msg.Message) bool {
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
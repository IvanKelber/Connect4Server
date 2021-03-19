package game

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"kelber.com/connect4/msg"
)

type RequestHandler struct {
	lobby          *Lobby
	activeSessions *map[string]*Session
	activeGames    *map[string]*Game
}

func CreateRequestHandler(lobby *Lobby,
	activeSessions *map[string]*Session,
	activeGames *map[string]*Game) RequestHandler {
	return RequestHandler{lobby, activeSessions, activeGames}
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
		fmt.Printf("Server does not handle responses.\n")
	default:
		fmt.Printf("Unknown message type: %d\n", message.Type)
	}
	return false

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
		rh.PlacePiece(session, message)

	case msg.ChallengePlayerReq:
		rh.CreatePlayerChallenge(session, message)

	case msg.ProposalAnswerReq:
		rh.ParseProposalAnswer(session, message)
	case msg.CancelProposalReq:

	default:
		fmt.Printf("Unknown request %d\n", message.ID)
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

func (rh *RequestHandler) CreatePlayerChallenge(session *Session, message msg.Message) {
	opponentUsername := string(message.Content[0])
	opponentSession := rh.lobby.GetSession(opponentUsername)
	opponentSession.SendChallengeProposal(session.GetUsername())
	session.WaitForChallengeResponse(opponentUsername)
}

func (rh *RequestHandler) ParseProposalAnswer(session *Session, message msg.Message) {
	accepted := message.Content[0][0] == msg.TrueByte
	if accepted {
		// Proposal was accepted
		username1 := string(message.Content[1])
		username2 := string(message.Content[2])

		player1 := rh.lobby.GetSession(username1)
		player2 := rh.lobby.GetSession(username2)
		if player1 != nil && player2 != nil {
			fmt.Printf("Starting game between %s and %s\n", username1, username2)
			id := fmt.Sprintf("%s%d%s", username1, time.Now().UnixNano(), username2)
			fmt.Println("Creating id: ", id)
			game := CreateGame(player1, player2, id)
			(*rh.activeGames)[id] = game
		}

	} else {
		// Proposal was rejected
		player := string(message.Content[1])
		rh.lobby.GetSession(player).ChallengeRejected()
		fmt.Printf("%s rejected the proposal from %s\n", session.GetUsername(), player)
	}
}

func (rh *RequestHandler) StartTurn(player int) {
	fmt.Printf("Player %d is starting their turn...\n", player)
}

func (rh *RequestHandler) PlacePiece(session *Session, message msg.Message) {
	gameId := string(message.Content[0])
	column, err := strconv.Atoi(string(message.Content[1]))
	if err == nil {
		if game, ok := (*rh.activeGames)[gameId]; ok {
			fmt.Println("Placing piece in column ", column)

			game.PlacePiece(session, column)
		}
	} else {
		fmt.Printf("Could not place piece because of error: %s\n", err)
	}

}

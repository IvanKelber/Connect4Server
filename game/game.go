package game

import (
	"fmt"
	"math/rand"
	"time"
)

//Constants for the game logic
const (
	PlayersPerGame = 2
	BoardRows      = 7
	BoardCols      = 6
)

//Game contains the state of a game: the players in the game, who's turn it is and
type Game struct {
	players           []*Session
	CurrentPlayerTurn byte
	board             [][]int
	Id                string
}

func emptyGame() Game {
	board := make([][]int, BoardRows)
	for i := range board {
		board[i] = make([]int, BoardCols)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Game{make([]*Session, 0), byte(r.Intn(2)), board, ""}
}

//CreateGame creates a game with two players
func CreateGame(player1, player2 *Session, id string) *Game {
	newGame := emptyGame()
	newGame.Id = id

	newGame.AddPlayer(player1)
	newGame.AddPlayer(player2)
	player1.StartGame(&newGame)
	player2.StartGame(&newGame)
	return &newGame
}

//AddPlayer adds a player to a game if possible
func (g *Game) AddPlayer(p *Session) {
	if !g.IsFull() {
		g.players = append(g.players, p)
	} else {
		fmt.Println("Failed to add player to full game")
	}
}

func (g Game) IsMyTurn(p *Session) bool {
	for i, player := range g.players {
		if p == player {
			return byte(i) == g.CurrentPlayerTurn
		}
	}
	return false
}

//IsFull allows a convenient method for determining if there is space
func (g Game) IsFull() bool {
	return len(g.players) == PlayersPerGame
}

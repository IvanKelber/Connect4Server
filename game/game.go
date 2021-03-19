package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//Constants for the game logic
const (
	PlayersPerGame = 2
	BoardRows      = 6
	BoardCols      = 7
)

//Game contains the state of a game: the players in the game, who's turn it is and
type Game struct {
	players           []*Session
	CurrentPlayerTurn int
	board             [][]int
	Id                string
}

func emptyGame() Game {
	board := make([][]int, BoardRows)
	for i := range board {
		board[i] = make([]int, BoardCols)
		for j := range board[i] {
			board[i][j] = -1 //Start off empty
		}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return Game{make([]*Session, 0), r.Intn(2), board, ""}
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

func (g *Game) PlacePiece(session *Session, column int) {
	playerId, err := g.GetPlayerId(session)
	if err == nil {
		g.board[g.findRowFromColumn(column)][column] = playerId
	} else {
		fmt.Printf("Could not place piece because of error: %s\n", err)
	}
	fmt.Println("Placed piece... board is now")
	fmt.Println(g)
}

func (g Game) findRowFromColumn(column int) int {
	i := BoardRows - 1
	for i >= 0 && g.board[i][column] >= 0 {
		i--
	}
	return i
}

func (g Game) GetPlayerId(p *Session) (int, error) {
	for i, player := range g.players {
		if p == player {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Failed to find id of player %s\n", p.GetUsername())
}

func (g Game) IsMyTurn(p *Session) bool {
	id, err := g.GetPlayerId(p)
	if err == nil {
		return id == g.CurrentPlayerTurn
	}
	fmt.Println(err)
	return false
}

//IsFull allows a convenient method for determining if there is space
func (g Game) IsFull() bool {
	return len(g.players) == PlayersPerGame
}

func (g Game) String() string {
	builder := strings.Builder{}
	for i := 0; i < len(g.board); i++ {
		builder.WriteString(fmt.Sprintf("%v\n", g.board[i]))
	}
	return builder.String()
}

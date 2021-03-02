package game

import "fmt"

//Constants for the game logic
const (
	PlayersPerGame = 2
	BoardSize      = 7
)

//Game contains the state of a game: the players in the game, who's turn it is and
type Game struct {
	players     []Player
	currentTurn int
	board       [][]int
}

func emptyGame() Game {
	board := make([][]int, BoardSize)
	for i := range board {
		board[i] = make([]int, BoardSize)
	}
	return Game{make([]Player, 0), 0, board}
}

//CreateGame creates a game with two players
func CreateGame(player1, player2 Player) *Game {
	newGame := emptyGame()
	newGame.AddPlayer(player1)
	newGame.AddPlayer(player2)
	player1.SetGame(&newGame)
	player2.SetGame(&newGame)
	return &newGame
}

//AddPlayer adds a player to a game if possible
func (g *Game) AddPlayer(p Player) {
	if !g.IsFull() {
		g.players = append(g.players, p)
	} else {
		fmt.Println("Failed to add player to full game")
	}
}

//IsFull allows a convenient method for determining if there is space
func (g Game) IsFull() bool {
	return len(g.players) == PlayersPerGame
}

package game

//Player is a server side player struct only
type Player struct {
	Username string
	Playing  bool
	game     *Game
}

//CreatePlayer creates a player not in a game
func CreatePlayer(username string) *Player {
	return &Player{username, false, nil}
}

//SetGame of the player
func (player *Player) SetGame(g *Game) {
	player.game = g
}

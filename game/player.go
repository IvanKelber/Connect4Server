package game

//Player is a server side player struct only
type Player struct {
	Username string
	Playing  bool
}

//CreatePlayer creates a player not in a game
func CreatePlayer(username string) *Player {
	return &Player{username, false}
}

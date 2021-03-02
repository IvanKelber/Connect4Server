package game

import (
	"net"
)

//Player is a server side player struct only
type Player struct {
	Username string
	Conn     *net.Conn
	Playing  bool
	game     *Game
}

//CreatePlayer creates a player not in a game
func CreatePlayer(username string, conn *net.Conn) Player {
	return Player{username, conn, false, nil}
}

//SetGame of the player
func (player *Player) SetGame(g *Game) {
	player.game = g
}

func (player Player) GetAddress() string {
	return (*player.Conn).RemoteAddr().String()
}

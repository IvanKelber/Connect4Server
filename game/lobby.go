package game

import (
	"fmt"
	"strings"
	"sync"
)

type lobbyType map[string]Player

//Lobby contains all players not currently in games
type Lobby struct {
	data map[string]Player
	*sync.Mutex
}

func (lobby Lobby) String() string {
	b := strings.Builder{}
	b.WriteString("{")
	for _, player := range lobby.data {
		b.WriteString(fmt.Sprintf("(%s, %s), ", player.Username, player.GetAddress()))
	}
	b.WriteString("}")
	return b.String()
}

//AddPlayer to lobby
func (lobby *Lobby) AddPlayer(p Player) {
	lobby.Lock()
	defer lobby.Unlock()
	if _, ok := lobby.data[p.Username]; !ok {
		lobby.data[p.Username] = p
		fmt.Printf("Successfully added player '%s' at %s\n", p.Username, p.GetAddress())
	}
}

//CreateLobby once at runtime for manifest of players
func CreateLobby() Lobby {
	data := make(map[string]Player)
	return Lobby{data, &sync.Mutex{}}
}

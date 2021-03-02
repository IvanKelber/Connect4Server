package game

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"kelber.com/connect4/msg"
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

func (lobby *Lobby) UpdatePlayers() {
	lobby.Lock()
	defer lobby.Unlock()
	usernames := make([][]byte, 0)
	for _, player := range lobby.data {
		usernames = append(usernames, []byte(player.Username))
	}
	m := msg.CreateNewMessage(msg.Response, msg.UpdateStateResp, 29, usernames)
	b := bytes.Buffer{}
	msg.Serialize(m, &b)
	for _, player := range lobby.data {
		(*player.Conn).Write(b.Bytes())
	}
}

//CreateLobby once at runtime for manifest of players
func CreateLobby() Lobby {
	data := make(map[string]Player)
	return Lobby{data, &sync.Mutex{}}
}

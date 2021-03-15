package game

import (
	"fmt"
	"strings"
	"sync"

	"kelber.com/connect4/msg"
)

//Lobby contains all players not currently in games
type Lobby struct {
	data map[string]*Session
	*sync.Mutex
}

func (lobby Lobby) String() string {
	b := strings.Builder{}
	b.WriteString("{")
	for username, sess := range lobby.data {
		b.WriteString(fmt.Sprintf("(%s, %s), ", username, sess.GetAddress()))
	}
	b.WriteString("}")
	return b.String()
}

//AddPlayer to lobby
func (lobby *Lobby) AddPlayerSession(username string, session *Session) bool {
	lobby.Lock()
	defer lobby.Unlock()
	if _, ok := lobby.data[username]; !ok {
		lobby.data[username] = session
		fmt.Printf("Successfully added player '%s' at %s\n", username, session.GetAddress())
		if _, otherName := lobby.data[session.GetUsername()]; otherName {
			// Session already has a username so we must delete the old value
			fmt.Println("Deleting former username: ", session.GetUsername())
			delete(lobby.data, session.GetUsername())
		}
	} else {
		fmt.Printf("Username %s already exists, failed to add to lobby\n", username)
		return false
	}
	return true
}

func (lobby *Lobby) Remove(session *Session) {
	lobby.Lock()
	defer lobby.Unlock()
	delete(lobby.data, session.GetUsername())
}

//UpdatePlayers is to keep the players in the lobby up to date with other users in the lobby
func (lobby *Lobby) UpdatePlayers() {
	lobby.Lock()
	defer lobby.Unlock()
	usernames := make([][]byte, 0)
	for _, player := range lobby.data {
		usernames = append(usernames, []byte(player.GetUsername()))
	}
	fmt.Println("Creating message: ", usernames)
	m := msg.CreateNewMessage(msg.Response, msg.UpdateLobbyResp, msg.DefaultContentDelimiter, usernames)
	for _, session := range lobby.data {
		fmt.Printf("Sending message to %s with address %v\n", session.GetUsername(), session.GetAddress())
		session.SendMessage(m)
	}
	fmt.Println()
	fmt.Printf("Lobby: %s\n", lobby)
}

func (lobby *Lobby) GetSession(username string) *Session {
	if session, ok := lobby.data[username]; ok {
		return session
	}
	return nil
}

//CreateLobby once at runtime for manifest of players
func CreateLobby() Lobby {
	data := make(map[string]*Session)
	return Lobby{data, &sync.Mutex{}}
}

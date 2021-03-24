package msg

import (
	"bytes"
	"fmt"
)

//MessageType is either Request or Response
type MessageType byte

//MessageID determines the type of request and response sent
type MessageID byte

//Message is serialized and sent via TCP
type Message struct {
	Type      MessageType
	ID        MessageID
	delimiter byte
	Content   [][]byte
}

const FalseByte byte = 0
const TrueByte byte = 1

//FieldDelimiter is used to deliniate sections in the byte data stream
const FieldDelimiter byte = 28

//DefaultContentDelimiter is used to deliniate bits of content
const DefaultContentDelimiter byte = 29

//EndOfMessage is used to denote an end to a message
const EndOfMessage byte = 31

//Request and Response consts
const (
	Request MessageType = iota
	Response
)

//Requst types
const (
	NewPlayerReq MessageID = iota
	StartGameReq
	StartTurnReq
	PlacePieceReq
	UpdateStateReq
	ChallengePlayerReq //sent from first player to server
	ProposalAnswerReq  //sent from second player to server in response to proposal
	CancelProposalReq
)

//Response types
const (
	NewPlayerResp MessageID = iota
	StartGameResp
	StartTurnResp
	PlacePieceResp        // Should contain information about win/loss/tie
	UpdateLobbyResp       // list of strings
	ChallengeProposalResp // string
	WaitForChallengeResp  // void
	ChallengeRejectedResp // bool sent to player who started the challenge
	CancelProposalResp
	GameOverResp
)

//CreateNewMessage is a constructor for Message
func CreateNewMessage(t MessageType, id MessageID, delimiter byte, content [][]byte) Message {
	return Message{t, id, delimiter, content}
}

//Serialize messages into a standardized byte format
func Serialize(message Message, buffer *bytes.Buffer) error {

	//First serialize the message type
	buffer.WriteByte(byte(message.Type))
	buffer.WriteByte(FieldDelimiter)

	//Second serialize the messageID
	buffer.WriteByte(byte(message.ID))
	buffer.WriteByte(FieldDelimiter)

	//Third serialize the delimiter (if any) to be used when parsing the content
	buffer.WriteByte(message.delimiter)
	buffer.WriteByte(FieldDelimiter)

	//Last append the content
	for _, content := range message.Content {
		_, err := buffer.Write(content)
		if err != nil {
			fmt.Printf("Failed to serialize %s\n", content)
			return err
		}
		buffer.WriteByte(message.delimiter)
	}
	buffer.WriteByte(FieldDelimiter)
	buffer.WriteByte(EndOfMessage)

	return nil
}

//Deserialize byte arrays into message given expected format
func Deserialize(buffer bytes.Buffer) Message {
	messageType, err := buffer.ReadBytes(FieldDelimiter)
	if err != nil {
		fmt.Printf("Failed to read message type %s\n", err)
	}
	messageType = messageType[0 : len(messageType)-1]

	messageID, err := buffer.ReadBytes(FieldDelimiter)
	if err != nil {
		fmt.Printf("Failed to read message ID %s\n", err)
	}
	messageID = messageID[0 : len(messageID)-1]

	contentDelimiter, _ := buffer.ReadByte()
	buffer.ReadByte() //Consumes last field delimiter before content begins

	content := make([][]byte, 0)

	for {
		c, err := buffer.ReadBytes(contentDelimiter)
		content = append(content, c[:len(c)-1])
		if err != nil {
			break
		}
	}
	return CreateNewMessage(MessageType(messageType[0]), MessageID(messageID[0]), contentDelimiter, content)
}

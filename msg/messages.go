package msg

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

//MessageType is either Request or Response
type MessageType int

//MessageID determines the type of request and response sent
type MessageID int

//Message is serialized and sent via TCP
type Message struct {
	Type      MessageType
	ID        MessageID
	delimiter byte
	Content   [][]byte
}

//DELIMITER is used to deliniate sections in the byte data stream
const DELIMITER byte = 28

//Request and Response consts
const (
	Request MessageType = iota
	Response
)

//Requst/Response types
const (
	NewPlayerReq MessageID = iota
	StartGameReq
	StartTurnReq
	PlacePieceReq
	UpdateStateReq
	AnimationDoneReq
	GameOverReq

	NewPlayerResp
	StartGameResp
	StartTurnResp
	PlacePieceResp
	UpdateStateResp
	AnimationDoneResp
	GameOverResp
)

//CreateNewMessage is a constructor for Message
func CreateNewMessage(t MessageType, id MessageID, delimiter byte, content [][]byte) Message {
	return Message{t, id, delimiter, content}
}

//Serialize messages into a standardized byte format
func Serialize(message Message, buffer *bytes.Buffer) error {

	//First serialize the message type
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(message.Type))
	buffer.Write(b)
	buffer.WriteByte(DELIMITER)
	//Second serialize the messageID
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(message.ID))
	buffer.Write(b)
	buffer.WriteByte(DELIMITER)

	//Third serialize the delimiter (if any) to be used when parsing the content
	buffer.WriteByte(message.delimiter)
	buffer.WriteByte(DELIMITER)
	//Last append the content
	for _, content := range message.Content {
		_, err := buffer.Write(content)
		if err != nil {
			fmt.Printf("Failed to serialize %s\n", content)
			return err
		}
		buffer.WriteByte(message.delimiter)
	}
	buffer.WriteByte(DELIMITER)
	buffer.WriteByte('\n')

	return nil
}

//Deserialize byte arrays into message given expected format
func Deserialize(buffer bytes.Buffer) Message {
	messageType, err := buffer.ReadBytes(DELIMITER)
	if err != nil {
		fmt.Printf("Failed to read message type %s\n", err)
	}
	messageType = messageType[0 : len(messageType)-1]

	messageID, err := buffer.ReadBytes(DELIMITER)
	if err != nil {
		fmt.Printf("Failed to read message ID %s\n", err)
	}
	messageID = messageID[0 : len(messageID)-1]

	contentDelimiter, _ := buffer.ReadByte()
	content := make([][]byte, 0)

	for {
		c, err := buffer.ReadBytes(contentDelimiter)
		content = append(content, c[:len(c)-1])
		if err != nil {
			break
		}
	}
	return CreateNewMessage(MessageType(binary.BigEndian.Uint16(messageType)), MessageID(binary.BigEndian.Uint16(messageID)), contentDelimiter, content)
}

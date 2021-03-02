package msg

import (
    "encoding/binary"
    "bytes"
	"fmt"
)

//MessageType is either Request or Response
type MessageType int

//MessageID determines the type of request and response sent
type MessageID int

//Message is serialized and sent via TCP
type Message struct {
	Type MessageType
	ID MessageID
	Content []byte
}

//DELIMITER is used to deliniate sections in the byte data stream
const DELIMITER byte = ' '

//Request and Response consts
const (
	Request MessageType = iota
	Response 
)

//Requst/Response types
const (
	StartGameReq MessageID = iota
	StartTurnReq
	PlacePieceReq
	UpdateStateReq
	AnimationDoneReq
	GameOverReq
	StartGameResp 
	StartTurnResp
	PlacePieceResp
	UpdateStateResp
	AnimationDoneResp
	GameOverResp
)

//CreateNewMessage is a constructor for Message
func CreateNewMessage(t MessageType, id MessageID, content []byte) Message {
	return Message{t, id, content}
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
	//Last append the content
	_ , err := buffer.Write(message.Content)
	if err != nil {
		fmt.Printf("Failed to serialize %s\n", message.Content)
		return err
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
	messageType = messageType[0:len(messageType) - 1]

	messageID, err := buffer.ReadBytes(DELIMITER)
	if err != nil {
		fmt.Printf("Failed to read message ID %s\n",err)
	}	
	messageID = messageID[0:len(messageID) - 1]

	content, err := buffer.ReadBytes(DELIMITER)
	if err != nil {
		fmt.Printf("Failed to read message content %s\n", err)
	}
	content = content[0:len(content) - 1]

	return CreateNewMessage(MessageType(binary.BigEndian.Uint16(messageType)), MessageID(binary.BigEndian.Uint16(messageID)), content)
}
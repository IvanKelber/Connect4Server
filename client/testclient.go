package main

import (
    "bufio"
    "fmt"
    "bytes"
    "net"
    "os"
    "kelber.com/connect4/msg"
)

const (
    connHost = "localhost"
    connPort = "8080"
    connType = "tcp"
)

func main() {
    fmt.Println("Connecting to " + connType + " server " + connHost + ":" + connPort)

    conn, err := net.Dial(connType, connHost+":"+connPort)
    if err != nil {
        fmt.Println("Error connecting:", err.Error())
        os.Exit(1)
    }
    reader := bufio.NewReader(os.Stdin)

    for {
        reader.ReadString('\n')
        
        m := msg.CreateNewMessage(msg.Request, msg.StartTurnReq, []byte{1})
        buffer := bytes.Buffer{}
        msg.Serialize(m, &buffer)
        fmt.Printf("msg: %v\n", buffer.Bytes())
        conn.Write(buffer.Bytes())
        
    }
}


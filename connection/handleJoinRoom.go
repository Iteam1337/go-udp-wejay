package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

func (c *Connection) handleJoinRoom() {
	msg := c.msg.(*message.JoinRoom)
	log.Println("handleJoinRoom", msg)
}

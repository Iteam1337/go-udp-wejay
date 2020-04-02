package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

func (c *Connection) handleLeaveRoom() {
	msg := c.msg.(*message.LeaveRoom)
	log.Println("handleLeaveRoom", msg)
}

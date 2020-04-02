package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
)

func (c *Connection) handleQueryRooms() {
	msg := c.msg.(*message.QueryRooms)
	log.Println("handleQueryRooms", msg)
}

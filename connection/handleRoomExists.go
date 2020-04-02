package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/rooms"
)

func (c *Connection) handleRoomExists() {
	msg := c.msg.(*message.RoomExists)
	log.Println("handleRoomExists", msg)

	res := message.RoomExistsResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	res.RoomId = msg.RoomId
	res.Ok = true
	res.Exists = rooms.Exists(msg.RoomId)

	c.send(&res)
}

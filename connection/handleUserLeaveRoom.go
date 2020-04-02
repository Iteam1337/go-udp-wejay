package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/rooms"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleUserLeaveRoom() {
	msg := c.msg.(*message.UserLeaveRoom)
	res := message.UserLeaveRoomResponse{Ok: false}
	log.Println("handleUserLeaveRoom", msg)

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	if exists := users.Exists(msg.UserId); !exists {
		res.Error = "user does not exist"
		c.send(&res)
		return
	}

	if ok := rooms.Evict(msg.UserId); ok {
		res.UserId = msg.UserId
		res.Ok = true
	}

	c.send(&res)
}

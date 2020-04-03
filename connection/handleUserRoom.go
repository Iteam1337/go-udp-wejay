package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleUserRoom() {
	msg := c.msg.(*message.UserRoom)
	res := message.UserRoomResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	if !users.Exists(msg.UserId) {
		res.Error = "user does not exist"
		c.send(&res)
		return
	}

	if user, err := users.GetUser(msg.UserId); err != nil {
		res.Error = err.Error()
		c.send(&res)
		return
	} else {
		res.UserId = msg.UserId
		res.RoomId = user.Room
		res.Ok = true
	}

	c.send(&res)
}

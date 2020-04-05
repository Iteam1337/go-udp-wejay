package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleUserExists() {
	msg := c.msg.(*message.UserExists)
	res := message.UserExistsResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	res.UserId = msg.UserId
	res.Ok = true
	res.Exists = users.Exists(msg.UserId)

	c.send(&res)
}

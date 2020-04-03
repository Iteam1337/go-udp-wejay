package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleDeleteUser() {
	msg := c.msg.(*message.DeleteUser)
	res := message.DeleteUserResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	id := msg.UserId
	res.UserId = id

	if exists := users.Exists(id); !exists {
		res.Error = "could not find user"
		c.send(&res)
		return
	}

	users.Delete(id)

	res.Ok = true
	c.send(&res)
}

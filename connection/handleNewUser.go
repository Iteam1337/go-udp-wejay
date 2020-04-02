package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleNewUser() {
	msg := c.msg.(*message.NewUser)
	log.Println("handleNewUser")
	res := message.NewUserResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	id := msg.UserId
	res.UserId = id
	if exists := users.Exists(id); exists {
		res.Error = "user exists"
		c.send(&res)
		return
	}

	users.New(id, msg.Code)
	res.Ok = true
	c.send(&res)
}

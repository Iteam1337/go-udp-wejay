package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleAction() {
	msg := c.msg.(*message.Action)
	log.Println("handleAction", msg)

	res := message.ActionResponse{Ok: false}

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

	user, err := users.GetUser(id)
	if err != nil {
		res.Error = err.Error()
		c.send(&res)
		return
	}

	if err := user.RunAction(msg.Action); err != nil {
		res.Error = err.Error()
		c.send(&res)
		return
	}

	res.Ok = true
	c.send(&res)
}

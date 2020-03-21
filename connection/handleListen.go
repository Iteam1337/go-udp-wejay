package connection

import (
	"fmt"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleListen() {
	var listen = make(chan user.ListenMsg, 1)
	msg := c.msg.(*message.Listen)
	fmt.Println("handleListen", msg)
	if msg == nil {
		return
	}

	id := msg.UserId
	if exists := users.Exists(id); !exists {
		return
	}

	user, err := users.GetUser(id)
	if err != nil {
		return
	}

	user.SetListen(&listen)

	for {
		msg := <-listen
		res := message.ListenResponse{
			UserId: id,
			Change: msg.Type,
			Meta:   msg.Meta,
			Ok:     true,
		}
		c.send(&res)
	}
}

package connection

import (
	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/spotifyauth"
)

func (c *Connection) handleCallbackURL() {
	msg := c.msg.(*message.CallbackURL)
	res := message.CallbackURLResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse message"
		c.send(&res)
		return
	}

	res.UserId = msg.UserId
	res.Ok = true
	res.Url = spotifyauth.AuthURL(msg.UserId)
	c.send(&res)
}

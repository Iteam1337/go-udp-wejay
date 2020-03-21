package connection

import (
	"log"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-udp-wejay/users"
)

func (c *Connection) handleNowPlaying() {
	msg := c.msg.(*message.NowPlaying)
	log.Println("handleNowPlaying", msg)

	res := message.NowPlayingResponse{Ok: false}

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

	var track message.Track
	if track, err = user.NowPlaying(); err != nil {
		res.Error = err.Error()
		c.send(&res)
		return
	}

	res.Ok = true
	res.Track = &track

	c.send(&res)
}

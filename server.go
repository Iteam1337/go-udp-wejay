package main

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/message"
	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-udp-wejay/inputtype"
	"github.com/Iteam1337/go-udp-wejay/user"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/golang/protobuf/proto"
)

type connection struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	it   types.InputType
	msg  proto.Message
}

func (c *connection) handleNowPlaying() {
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

	if exists := user.Exists(id); !exists {
		res.Error = "could not find user"
		c.send(&res)
		return
	}

	user, err := user.GetUser(id)
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

func (c *connection) handleAction() {
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

	if exists := user.Exists(id); !exists {
		res.Error = "could not find user"
		c.send(&res)
		return
	}

	user, err := user.GetUser(id)
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

func (c *connection) handleUserExists() {
	msg := c.msg.(*message.UserExists)
	log.Println("handleUserExists", msg)

	res := message.UserExistsResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	res.UserId = msg.UserId
	res.Ok = true
	res.Exists = user.Exists(msg.UserId)

	c.send(&res)
}

func (c *connection) handleNewUser() {
	msg := c.msg.(*message.NewUser)
	log.Println("handleNewUser", msg)
	res := message.NewUserResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse input"
		c.send(&res)
		return
	}

	id := msg.UserId
	res.UserId = id
	if exists := user.Exists(id); exists {
		res.Error = "user exists"
		c.send(&res)
		return
	}

	user.NewUser(id, msg.Code)
	res.Ok = true
	c.send(&res)
}

func (c *connection) handleCallbackURL() {
	msg := c.msg.(*message.CallbackURL)
	log.Println("handleCallbackURL", msg)
	res := message.CallbackURLResponse{Ok: false}

	if msg == nil {
		res.Error = "could not parse message"
		c.send(&res)
		return
	}

	res.UserId = msg.UserId
	res.Ok = true
	res.Url = user.AuthURL(msg.UserId)
	c.send(&res)

}

func (c *connection) handlePing() {
	utils.SendM(c.it.Inv(), &message.Pong{Int: 0}, c.conn, c.addr)
}

func (c *connection) send(r proto.Message) {
	utils.SendM(c.it.Inv(), r, c.conn, c.addr)
}

func (c connection) parse(conn *net.UDPConn) (addr *net.UDPAddr, it types.InputType, buffer []byte, err error) {
	var length int
	buffer = make([]byte, 4096)
	length, addr, err = conn.ReadFromUDP(buffer)
	if err != nil {
		log.Println(err)
		return
	}
	it, err = inputtype.FromBuffer(buffer[:2])
	if err != nil {
		log.Println(err)
		utils.SendEmpty(conn, addr)
		return
	}
	buffer = buffer[2:length]
	return
}

func (c connection) read(msg proto.Message, buf []byte) (err error) {
	if err = proto.Unmarshal(buf[:], msg); err != nil {
		log.Println(err)
		return
	}
	return
}

func (c connection) convert(it types.InputType, buf []byte) (pb proto.Message) {
	pb = it.Message()
	c.read(pb, buf)
	return
}

func parseConn(conn *net.UDPConn) (c connection, e error) {
	addr, it, buffer, err := c.parse(conn)
	if err != nil {
		e = err
		return
	}

	c.conn = conn
	c.addr = addr
	c.it = it
	c.msg = c.convert(it, buffer)
	return
}

func listen(address string) {
	newConnection := func(c *net.UDPConn) {
		var con connection
		var err error
		if con, err = parseConn(c); err != nil {
			return
		}

		switch con.it {
		case types.IAction:
			con.handleAction()
		case types.IUserExists:
			con.handleUserExists()
		case types.INewUser:
			con.handleNewUser()
		case types.ICallbackURL:
			con.handleCallbackURL()
		case types.IPing:
			con.handlePing()
		case types.INowPlaying:
			con.handleNowPlaying()
		default:
			utils.SendEmpty(con.conn, con.addr)
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	log.Printf("listening on UDP %s\n", address)
	for {
		newConnection(ln)
	}
}

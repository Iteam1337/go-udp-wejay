package connection

import (
	"log"
	"net"

	"github.com/Iteam1337/go-protobuf-wejay/types"
	"github.com/Iteam1337/go-udp-wejay/inputtype"
	"github.com/Iteam1337/go-udp-wejay/utils"
	"github.com/golang/protobuf/proto"
)

// Connection …
type Connection struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	it   types.InputType
	msg  proto.Message
}

func (c *Connection) send(r proto.Message) {
	utils.SendM(c.it.Inv(), r, c.conn, c.addr)
}

func (c Connection) parse(conn *net.UDPConn) (addr *net.UDPAddr, it types.InputType, buffer []byte, err error) {
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

func (c Connection) read(msg proto.Message, buf []byte) (err error) {
	if err = proto.Unmarshal(buf[:], msg); err != nil {
		log.Println(err)
		return
	}
	return
}

func (c Connection) convert(it types.InputType, buf []byte) (pb proto.Message) {
	pb = it.Message()
	c.read(pb, buf)
	return
}

// Handler …
func (c *Connection) Handler() {
	go func() {
		switch c.it {
		case types.IAction:
			c.handleAction()
		case types.IUserExists:
			c.handleUserExists()
		case types.INewUser:
			c.handleNewUser()
		case types.ICallbackURL:
			c.handleCallbackURL()
		case types.IPing:
			c.handlePing()
		case types.INowPlaying:
			c.handleNowPlaying()
		case types.IListen:
			c.handleListen()
		case types.IDeleteUser:
			c.handleDeleteUser()
		default:
			utils.SendEmpty(c.conn, c.addr)
		}
	}()
}

// ParseConn …
func ParseConn(conn *net.UDPConn) (c Connection, e error) {
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

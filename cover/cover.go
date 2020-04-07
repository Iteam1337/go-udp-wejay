package cover

import (
	"encoding/binary"
	"log"
	"net"

	"github.com/Iteam1337/go-udp-wejay/utils"
)

var (
	source = utils.GetEnv("GEN_COVER", "localhost:8091")
)

func Gen(label string) (buf []byte, err error) {
	var (
		conn *net.TCPConn
		addr *net.TCPAddr
	)

	log.Println("new request for", source)

	if addr, err = net.ResolveTCPAddr("tcp", source); err != nil {
		log.Println("cant resolve", err)
		return
	}

	if conn, err = net.DialTCP("tcp", nil, addr); err != nil {
		log.Println("cant connect", err)
		return
	}

	defer conn.Close()

	if _, err = conn.Write([]byte(label)); err != nil {
		log.Println("cant send", err)
		return
	}

	bufferSize := make([]byte, 4)
	if _, err = conn.Read(bufferSize); err != nil {
		log.Println("cant read buffer size", err)
		return
	}

	bufferSizeParsed := binary.LittleEndian.Uint32(bufferSize)
	msg := make([]byte, bufferSizeParsed)
	if _, err = conn.Read(msg); err != nil {
		log.Println("cant read", err)
		return
	}

	buf = msg

	return
}

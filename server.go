package main

import (
	"log"
	"net"

	"github.com/Iteam1337/go-udp-wejay/connection"
)

func listen(address string) {
	newConnection := func(c *net.UDPConn) {
		if con, err := connection.ParseConn(c); err == nil {
			con.Handler()
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

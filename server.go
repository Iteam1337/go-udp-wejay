package main

import (
	"log"
	"net"
)

// Listen …
func Listen(address string) {
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

package main

import (
	"fmt"
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

	fmt.Printf("listening on UDP %s\n", address)
	for {
		NewConnection(ln)
	}
}

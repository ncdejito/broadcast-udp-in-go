package main

import (
	"net"
)

func main() {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:8829")
	if err != nil {
		panic(err)
	}

	_, err2 := pc.WriteTo([]byte("data to transmit"), addr)
	if err2 != nil {
		panic(err2)
	}
}

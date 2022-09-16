package main

import (
	"fmt"
	"net"
)

func main() {
	address := "127.0.0.1"
	port := 1234

	// listenAddr := ":0"
	// listenAddr = fmt.Sprintf("%s:%d", address, port)

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	pc, err := net.ListenUDP("udp4", &addr)
	// pc, err := net.ListenPacket("udp4", listenAddr)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Printf("Listening on %s\n", pc.LocalAddr().String())
	buf := make([]byte, 1024)
	n, senderAddr, err := pc.ReadFromUDP(buf)
	// n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s sent this: %s\n", senderAddr, buf[:n])

	response := fmt.Sprintf("From server: got this: %s", buf[:n])
	_, err2 := pc.WriteToUDP([]byte(response), senderAddr)
	if err2 != nil {
		panic(err2)
	}
}

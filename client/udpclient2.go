package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	address := "127.0.0.1"
	port := 1234

	listenAddr := ":0"
	listenAddr = fmt.Sprintf("%s:%d", address, port)

	pc, err := net.Dial("udp4", listenAddr)
	// pc, err := net.ListenPacket("udp4", listenAddr)
	if err != nil {
		panic(err)
	}

	// these 2 addresses output different strings for some reason
	fmt.Printf("Listening on %s\n", listenAddr)
	fmt.Printf("Listening on %s\n", pc.LocalAddr().String())
	buf := make([]byte, 1024)

	// addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:12345")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = pc.WriteTo([]byte("data to transmit"), addr)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Fprintf(pc, "data to transmit")
	_, err = bufio.NewReader(pc).Read(buf)
	fmt.Printf("%s\n", buf)

	defer pc.Close()
}

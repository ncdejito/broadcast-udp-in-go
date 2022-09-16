package main

import (
	"bufio"
	"fmt"
	"net"
)

// https://stackoverflow.com/questions/26028700/write-to-client-udp-socket-in-go

func main() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp4", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

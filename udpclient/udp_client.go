package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	add := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 5683,
	}
	conn, err := net.DialUDP("udp", nil, add)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(3 * time.Second)
		conn.Write([]byte("hello udp tcpserver"))
	}
}

package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	ip := net.ParseIP("224.0.0.250")
	srcAddr := &net.UDPAddr{IP: []byte{127, 0, 0, 1}, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9981}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		time.Sleep(3 * time.Second)
		_, err := conn.Write([]byte("hello"))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("<%s>\n", conn.RemoteAddr())
	}
}

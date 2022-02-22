package main

import (
	"fmt"
	"net"
	toy_server "toy-server"
)

func main() {
	udpServer := &toy_server.Server{
		UdpProtocolName: "udp4",
		UdpBindHost:     "127.0.0.1",
		UdpPort:         5683,
		UdpHandler: func(conn net.UDPConn) error {
			for {
				buf := make([]byte, 1024)
				size, _ := conn.Read(buf)
				fmt.Println(string(buf[:size]))
			}
		},
	}

	udpServer.Serve()
}

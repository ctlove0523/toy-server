package main

import (
	"fmt"
	"net"
	toy_server "toy-server"
)

func main() {
	mserver := &toy_server.Server{
		MulticastBindHost: "224.0.0.250",
		MulticastPort:     9981,
		MulticastHandler: func(conn net.UDPConn) error {
			for {
				buf := make([]byte, 1024)
				size, addr, err := conn.ReadFromUDP(buf)
				if err != nil {
					fmt.Println(err)
					return err
				}
				fmt.Println(addr)
				fmt.Println(string(buf[:size]))
			}
		},
	}
	mserver.Serve()
}

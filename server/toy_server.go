package main

import (
	"fmt"
	"io/ioutil"
	"net"
	toy_server "toy-server"
)

func main() {
	ca, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		fmt.Printf("load ca failed,err = %s\n", err)
		return
	}

	serverCert, err := ioutil.ReadFile("./certs/server.pem")
	if err != nil {
		fmt.Printf("load server cert failed,err = %s\n", err)
		return
	}

	serverKey, err := ioutil.ReadFile("./certs/server.key")
	if err != nil {
		fmt.Printf("load server key failed,err = %s\n", err)
		return
	}

	server := &toy_server.Server{
		TcpProtocolName:      "tcp",
		TcpBindHost:          "127.0.0.1",
		TcpPort:              3883,
		TlsEnabled:           true,
		CaCertificate:        ca,
		ServerCertificate:    serverCert,
		ServerCertificateKey: serverKey,
		TcpHandler: func(conn net.Conn) error {
			for {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Println(err)
					return err
				}
				fmt.Println(string(buf[:n]))
			}
		},
	}

	server.Serve()
}

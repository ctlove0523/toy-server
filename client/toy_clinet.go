package main

import (
	"crypto/tls"
	"fmt"
)

func main() {
	cipherSuites := []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}
	config:=&tls.Config{
		MaxVersion:         tls.VersionTLS12,
		MinVersion:         tls.VersionTLS12,
		CipherSuites:       cipherSuites,
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:3883", config)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("connect to server success")
	conn.Write([]byte("hello tls server"))
}

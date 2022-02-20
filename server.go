package toy_server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

type TcpConnHandler func(conn net.Conn) error
type UdpConnHandler func(conn net.UDPConn) error

type Server struct {
	ProtocolName         string
	BindHost             string
	Port                 uint16
	TlsEnabled           bool
	CaCertificate        []byte
	ServerCertificate    []byte
	ServerCertificateKey []byte
	Handler              TcpConnHandler
}

func (s *Server) createTcpListener() net.Listener {
	// process tcp protocol
	if s.TlsEnabled {
		// 处理ca证书
		serverCaPool := x509.NewCertPool()
		serverCaPool.AppendCertsFromPEM(s.CaCertificate)

		// 处理服务端证书
		serverCert, err := tls.X509KeyPair(s.ServerCertificate, s.ServerCertificateKey)
		if err != nil {
			fmt.Println("load server certificate failed")
			return nil
		}
		var serverCerts []tls.Certificate
		serverCerts = append(serverCerts, serverCert)

		// 加密套件
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

		config := &tls.Config{
			RootCAs:      serverCaPool,
			Certificates: serverCerts,
			CipherSuites: cipherSuites,
			MaxVersion:   tls.VersionTLS12,
			MinVersion:   tls.VersionTLS12,
		}

		listener, err := tls.Listen(s.ProtocolName, fmt.Sprintf("%s:%d", s.BindHost, s.Port), config)
		if err != nil {
			fmt.Printf("tls server bind %s,port %d failed,reason %s\n", s.BindHost, s.Port, err)
			return nil
		}

		return listener
	}
	listener, err := net.Listen(s.ProtocolName, fmt.Sprintf("%s:%d", s.BindHost, s.Port))
	if err != nil {
		fmt.Printf("tcp server bind %s,port %d failed,reason %s", s.BindHost, s.Port, err)
		return nil
	}

	return listener

}

func (s *Server) Serve() {
	var listener net.Listener

	listener = s.createTcpListener()
	if listener == nil {
		return
	}

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("server accept new connectin faield %s\n", err)
			continue
		}

		go func() {
			err := s.Handler(conn)
			if err != nil {
				fmt.Printf("begin to close connection %s\n", conn.RemoteAddr())
				err = conn.Close()
				if err != nil {
					fmt.Printf("close connection failed %s\n", err)
				}

				return
			}
		}()
	}

}

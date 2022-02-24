package toy_server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
)

type TcpConnHandler func(conn net.Conn) error
type UdpConnHandler func(conn net.UDPConn) error
type MulticastHandler func(conn net.UDPConn) error

type Server struct {
	TcpProtocolName      string
	TcpBindHost          string
	TcpPort              uint16
	TlsEnabled           bool
	CaCertificate        []byte
	ServerCertificate    []byte
	ServerCertificateKey []byte
	TcpHandler           TcpConnHandler
	UdpProtocolName      string
	UdpBindHost          string
	UdpPort              int
	UdpHandler           UdpConnHandler
	MulticastBindHost    string
	MulticastPort        int
	MulticastHandler     MulticastHandler
	stateFlag            chan struct{}
}

func NewTcpServer(host string, port uint16, handler TcpConnHandler) *Server {
	if handler == nil {
		fmt.Println("tcp server must set one TcpConnHandler")
		return nil
	}

	s := &Server{
		TcpProtocolName: "tcp",
		TcpBindHost:     host,
		TcpPort:         port,
		TcpHandler:      handler,
		TlsEnabled:      false,
		stateFlag:       make(chan struct{}),
	}

	return s
}

func NewTlsServer(host string, port uint16, handler TcpConnHandler, ca, serverCert, serverKey []byte) *Server {
	if handler == nil {
		fmt.Println("tls server must set one TcpConnHandler")
		return nil
	}

	s := &Server{
		TcpProtocolName:      "tcp",
		TcpBindHost:          host,
		TcpPort:              port,
		TlsEnabled:           true,
		CaCertificate:        ca,
		ServerCertificate:    serverCert,
		ServerCertificateKey: serverKey,
		TcpHandler:           handler,
		stateFlag:            make(chan struct{}),
	}

	return s
}

func (s *Server) Serve() {
	// 检查是否开启UDP服务
	if len(s.UdpProtocolName) != 0 {
		fmt.Printf("begin to serve udp,listen = %s,port = %d\n", s.UdpBindHost, s.UdpPort)

		udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.UdpBindHost, s.UdpPort))
		if err != nil {
			fmt.Printf("resovle udp addr failed,error = %s\n", err)
			return
		}

		udpConn, err := net.ListenUDP(s.UdpProtocolName, udpAddr)
		if err != nil {
			fmt.Println("serve udp failed")
			return
		}
		go func() {
			err := s.UdpHandler(*udpConn)
			if err != nil {
				err = udpConn.Close()
				if err != nil {
					fmt.Println("close udp connection failed")
				}
			}
		}()
	}

	// 检查是否多播
	if len(s.MulticastBindHost) != 0 {
		fmt.Println(fmt.Sprintf("%s:%d", s.MulticastBindHost, s.MulticastPort))
		gaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s.MulticastBindHost, s.MulticastPort))
		if err != nil {
			fmt.Printf("resovle multicast udp addr failed,error = %s\n", err)
			return
		}
		conn, err := net.ListenMulticastUDP("udp", nil, gaddr)
		if err != nil {
			fmt.Printf("listen multicast udp failed,error = %s\n", err)
			return
		}

		go func() {
			for {
				err := s.MulticastHandler(*conn)
				if err != nil {
					fmt.Println("process multicast message failed")
					err = conn.Close()
					if err != nil {
						fmt.Println("close connection failed")
					}
					return
				}

			}
		}()
	}

	if len(s.TcpProtocolName) != 0 {
		var listener net.Listener

		listener = s.createTcpListener()
		if listener == nil {
			return
		}

		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("tcpserver accept new connectin faield %s\n", err)
				continue
			}

			go func() {
				err := s.TcpHandler(conn)
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

	<-s.stateFlag
}

func (s *Server) lookUpHost(host string) net.IP {
	ip, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("can't resolve host %s\n", host)
		return nil
	}

	if len(ip) == 0 {
		fmt.Printf("host resolved but no ip address %s\n", host)
		return nil
	}

	return ip[0]

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
			fmt.Println("load tcpserver certificate failed")
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

		listener, err := tls.Listen(s.TcpProtocolName, fmt.Sprintf("%s:%d", s.TcpBindHost, s.TcpPort), config)
		if err != nil {
			fmt.Printf("tls tcpserver bind %s,port %d failed,reason %s\n", s.TcpBindHost, s.TcpPort, err)
			return nil
		}

		return listener
	}
	listener, err := net.Listen(s.TcpProtocolName, fmt.Sprintf("%s:%d", s.TcpBindHost, s.TcpPort))
	if err != nil {
		fmt.Printf("tcp tcpserver bind %s,port %d failed,reason %s", s.TcpBindHost, s.TcpPort, err)
		return nil
	}

	return listener

}

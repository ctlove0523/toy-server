# toy server

A simple tcp and udp server



### 1 How to create tcp server

**method 1:**

~~~go
s:=&Server{
	TcpProtocolName: "tcp",
	TcpBindHost:     "localhost",
	TcpPort:         5230,
	TcpHandler: func(conn net.Conn) error {
        // you code to process conn
		return nil
	},
}
s.Serve()
~~~

**method 2：**`NewTcpServer`

~~~go
handler := func(conn net.Conn) error {
	// your code process conn
	return nil
}
s := toy_server.NewTcpServer("localhost", 5230, handler)
s.Serve()
~~~


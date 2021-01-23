package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	tp "github.com/Asphaltt/go-tproxy"
)

var (
	flgListen = flag.String("l", ":443", "local address to listen on")
)

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *flgListen)
	if err != nil {
		fmt.Printf("failed to resolve listen address:%s, err:%v\n", *flgListen, err)
		return
	}

	lis, err := tp.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Printf("failed to listen on the address:%s, err:%v\n", *flgListen, err)
		return
	}
	defer lis.Close()

	fmt.Println("listen on address:", *flgListen)

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Printf("failed to accept new connection, err:%v\n", err)
			return
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client address is", conn.RemoteAddr().String())
	fmt.Println("server address is", conn.LocalAddr().String())

	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("failed to recv message from client, err:%v\n", err)
		return
	}

	fmt.Println("recv msg:", string(buf[:n]))

	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte("world")); err != nil {
		fmt.Printf("failed to send message to client, err:%v\n", err)
		return
	}

	fmt.Println("sent msg:", "world")
}

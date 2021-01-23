package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	tp "github.com/Asphaltt/go-tproxy"
)

var (
	flgLocal  = flag.String("l", "1.1.1.1:43562", "local address for client dial to server")
	flgRemote = flag.String("r", "2.2.2.2:8888", "remote address for client to connect")
)

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *flgLocal)
	if err != nil {
		fmt.Printf("failed to resolve local address: %s, err: %v\n", *flgLocal, err)
		return
	}

	raddr, err := net.ResolveTCPAddr("tcp", *flgRemote)
	if err != nil {
		fmt.Printf("failed to resolve remote address: %s, err: %v\n", *flgRemote, err)
		return
	}

	conn, err := tp.DialTCP(laddr, raddr)
	if err != nil {
		fmt.Printf("failed to dial to server, err: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("local  address is", laddr.String())
	fmt.Println("remote address is", raddr.String())

	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte("hello")); err != nil {
		fmt.Printf("failed to send hello, err: %v\n", err)
		return
	}

	fmt.Println("sent msg:", "hello")

	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("failed to recv message from server, err: %v\n", err)
		return
	}

	fmt.Println("recv msg:", string(buf[:n]))
}

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync"

	tproxy "github.com/Asphaltt/go-tproxy"
)

func main() {
	var listenAddr, listenInterface, forwardInterface string
	flag.StringVar(&listenAddr, "L", ":9999", "local address to listen on for tproxy")
	flag.StringVar(&listenInterface, "l", "", "an interface to bind with listener")
	flag.StringVar(&forwardInterface, "f", "", "an interface to forward data")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		fmt.Println("failed to resolve local addr", listenAddr, "err:", err)
		return
	}

	lis, err := tproxy.ListenTCPWithDevice(listenInterface, "tcp", laddr)
	if err != nil {
		fmt.Println("failed to listen on", listenAddr, "err:", err)
		return
	}
	defer lis.Close()

	for {
		conn, err := lis.(*tproxy.Listener).AcceptTProxy()
		if err != nil {
			fmt.Println("failed to accept new connection, err:", err)
			return
		}

		go proxyConn(conn, forwardInterface)
	}
}

func proxyConn(lconn *tproxy.Conn, device string) {
	defer lconn.Close()

	rconn, err := lconn.DialOriginalDestinationWithDevice(device, false)
	if err != nil {
		fmt.Println("failed to dial to", lconn.LocalAddr(), "err:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	ioCopy := func(src io.Reader, dst io.Writer) {
		io.Copy(dst, src)
		wg.Done()
	}
	go ioCopy(lconn, rconn)
	go ioCopy(rconn, lconn)
	wg.Wait()
}

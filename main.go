package main

import (
	"fmt"
	"net"
	"flag"
	"os"

	"github.com/Aspirin4k/chat-server/p2p"
	"github.com/Aspirin4k/chat-server/cui"
)

func main() {
	host := flag.String("host", "127.0.0.1", "ip-address")
	port := flag.Int("port", 7777, "port")
	remote := flag.String("remote", "", "remote address to join the network")
	flag.Parse()

	addr := net.ParseIP(*host)
	if addr == nil {
		fmt.Fprint(os.Stderr,"Invalid address")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout,"Using address %s:%d\n", addr.String(), *port)

	var addrRemote net.IP
	addrRemote = nil
	if *remote != "" {
		addrRemote = net.ParseIP(*remote)
		fmt.Fprintf(os.Stdout,"Remote address %s:%d\n", addrRemote.String(), *port)
	}

	go cui.Render()

	p2p.CreateAndListen(addr, addrRemote)
}
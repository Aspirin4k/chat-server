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
	flag.Parse()

	addr := net.ParseIP(*host)
	if addr == nil {
		fmt.Fprint(os.Stderr,"Invalid address")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout,"Using address %s:%d\n", addr.String(), *port)

	go cui.Render()

	p2p.CreateAndListen(addr)
}
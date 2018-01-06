package main

import (
	"flag"
	"net"
	"fmt"
	"os"

	"github.com/Aspirin4k/chat-server/p2p/network-operations"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
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

	// p2p.Create(addr, *port, *port)
	// network_operations.Join(p2p.ServerAddress, p2p.ServerID, p2p.ServerAddress)
	// network_operations.JoinAddBefore(p2p.ServerAddress, p2p.ServerID, "192.168.1.103")
	// network_operations.ReceiveIDs(p2p.ServerAddress, p2p.ServerID, p2p.ServerID)
	//network_operations.SendMessage(
	//	network_operations.ParseAddress(*remote, 7778),
	//	fmt.Sprintf("%d %s", declarations.CLIENT_HELLO, "asaasb"))
	network_operations.SendMessage(
		network_operations.ParseAddress(*remote, 7778),
		fmt.Sprintf("%d %s %s %d", declarations.CLIENT_REGISTER, "aazxd", "znt", 12345))
}
package network_operations

import (
	"net"
	"fmt"
	"os"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

func SendMessage(addr *net.TCPAddr, message string) {
	conn, err := net.DialTCP("tcp", nil, addr)
	error_catcher.CheckError(err)
	fmt.Fprint(os.Stdout,"Connected to remote node..\n")
	fmt.Fprintf(os.Stdout,"Sending: %s\n", message)

	_, err = conn.Write([]byte(message))
}

func ParseAddress(address string) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp4",
		fmt.Sprintf("%s:%d", address, declarations.PORT))
	error_catcher.CheckError(err)

	return tcpAddr
}
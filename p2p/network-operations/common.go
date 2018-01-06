package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/error-catcher"
)

func SendMessage(addr *net.TCPAddr, message string) {
	conn, err := net.DialTCP("tcp", nil, addr)
	error_catcher.CheckError(err)
	error_catcher.PushMessage(fmt.Sprintf("Sending: %s to %s", message, addr.IP.String()))

	_, err = conn.Write([]byte(message))
}

func ParseAddress(address string, port int) *net.TCPAddr {
	tcpAddr, err := net.ResolveTCPAddr("tcp4",
		fmt.Sprintf("%s:%d", address, port))
	error_catcher.CheckError(err)

	return tcpAddr
}
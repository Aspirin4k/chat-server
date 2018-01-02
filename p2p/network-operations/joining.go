package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

/**
Метод запроса подключения узла к существующей сети
 */
func Join(addr *net.TCPAddr, serverID int, serverIP *net.TCPAddr) {
	SendMessage(addr,
		fmt.Sprintf("%d %d %s", declarations.NODE_JOINING, serverID, serverIP.IP.String()))
}

/**
Метод запроса подключения данного узла перед целевым
 */
func JoinAddBefore(addr *net.TCPAddr, remoteID int, remoteNode string) {
	SendMessage(addr,
		fmt.Sprintf("%d %d %s", declarations.NODE_JOINING_ADD_BEFORE, remoteID, remoteNode))
}
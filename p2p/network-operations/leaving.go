package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

/**
Посылает сообщение о покидании сети
 */
func Leave(addr *net.TCPAddr, remoteId, serverId int) {
	// Если мы не единственный узел в сети
	if remoteId != serverId {
		msg := fmt.Sprintf("%d %d", declarations.NODE_LEAVING, serverId)
		msg += fmt.Sprintf("\n%s", PrepareResources())
		SendMessage(addr, msg)
	}
}

/**
Сообщает другим узлам о том, что данный покинул сеть
 */
func Leaved(addr *net.TCPAddr, serverAddress *net.TCPAddr, serverId int, remoteId int) {
	SendMessage(addr,
		fmt.Sprintf("%d %d\n%d %s", declarations.NODE_LEAVED, remoteId, serverId, serverAddress.IP.String()))
}
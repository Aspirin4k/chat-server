package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

/**
Сообщает узлу о том, что он должен обновить свою пальцевую таблицу
@nodes - массив из пар Идентификатор - Адрес
 */
func UpdateFingers(addr *net.TCPAddr, serverId int, serverAddress string, nodes [][]string) {
	msg := fmt.Sprintf("%d %d", declarations.FINGERS_UPDATE, serverId)
	for _, val := range nodes {
		msg += fmt.Sprintf("\n%s %s", val[0], val[1])
	}
	msg += fmt.Sprintf("\n%d %s", serverId, serverAddress)

	SendMessage(addr, msg)
}
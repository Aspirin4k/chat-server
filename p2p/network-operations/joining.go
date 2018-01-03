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

/**
Просим удаленный хост добавить нас в его пальцевую таблицу
 */
func AddMeToFinger(addr *net.TCPAddr, serverID int, serverIP string) {
	SendMessage(addr,
		fmt.Sprintf("%d %d %s", declarations.NODE_ADD_ME_TO_FINGER, serverID, serverIP))
}

/**
Просим удаленный хост добавить данный узел в его пальцевую таблицу
 */
func AddMeToFingerMessage(addr *net.TCPAddr, serverID int, serverIP string, message string) {
	SendMessage(addr,
		fmt.Sprintf("%s\n%d %s", message, serverID, serverIP))
}
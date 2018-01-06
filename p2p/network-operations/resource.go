package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
)

/**
Посылает часть таблицы рессурсов на удаленный хост.
@address - адрес, на который посылаются записи
remoteID - идентификатор удаленного хоста
serverID - идентификатор локального хоста
 */
func ReceiveIDs(address *net.TCPAddr, remoteID int, serverID int) {
	var message string
	message = fmt.Sprintf("%d %d", declarations.RESOURCE_RECEIVE_IDS, serverID)
	//for _, v := range tables.ResourcesIDsTable {
	//	if (v.ID < remoteID) && (v.ID > serverID) && (remoteID > serverID) ||
	//			(remoteID < serverID) && ((v.ID > serverID) || (v.ID < remoteID)) {
	//		message += fmt.Sprintf("\n%d %d %s", v.ID, v.HostID, v.Address.IP.String())
	//		tables.ResourceRemoveByKey(v.ID);
	//	}
	//}
	for _, v := range tables.ActiveClientsTable {
		if (v.ClientID < remoteID) && (v.ClientID > serverID) && (remoteID > serverID) ||
				(remoteID < serverID) && ((v.ClientID > serverID) || (v.ClientID < remoteID)) {
					message += fmt.Sprintf("\n%d %s", v.ClientID, v.Address.IP.String())
					tables.RemoveActiveClientByKey(v.ClientID)
		}
	}

	SendMessage(address, message)
}

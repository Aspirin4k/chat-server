package network_operations

import (
	"net"
	"fmt"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
)

func ReceiveIDs(address *net.TCPAddr, remoteID int, serverID int) {
	var message string
	message = fmt.Sprintf("%d\n", declarations.RESOURCE_RECEIVE_IDS)
	for _, v := range tables.ResourcesIDsTable {
		if (v.ID < remoteID) && (v.ID > serverID) && (remoteID > serverID) ||
				(remoteID < serverID) && ((v.ID > serverID) || (v.ID < remoteID)) {
			message += fmt.Sprintf("%d %s\n", v.ID, v.Address.IP.String())
			tables.ResourceRemoveByKey(v.ID);
		}
	}

	SendMessage(address, message)
}

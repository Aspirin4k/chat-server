package tables

import (
	"net"

	"github.com/sadlil/go-trigger"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

var ActiveClientsTable []declarations.ActiveClient

func AddActiveClient(addr *net.TCPAddr, id int) {
	ActiveClientsTable =
		append(ActiveClientsTable, declarations.ActiveClient{id,addr})
	AddAllActiveClient(addr, id)

	trigger.Fire(declarations.ACTIVE_CLIENTS_CHANGED, ActiveClientsTable)
}

func RemoveActiveClientByKey(id int) {
	for i, v := range ActiveClientsTable {
		if v.ClientID == id {
			ActiveClientsTable = append(ActiveClientsTable[:i], ActiveClientsTable[i+1:]...)
		}
	}

	trigger.Fire(declarations.ACTIVE_CLIENTS_CHANGED, ActiveClientsTable)
}

func RemoveActiveClientByIndex(i int) {
	ActiveClientsTable = append(ActiveClientsTable[:i], ActiveClientsTable[i+1:]...)

	trigger.Fire(declarations.ACTIVE_CLIENTS_CHANGED, ActiveClientsTable)
}
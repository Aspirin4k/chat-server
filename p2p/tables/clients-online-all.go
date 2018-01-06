package tables

import (
	"sync"
	"net"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/error-catcher"
	"fmt"
)

var AllActiveClients 		[]declarations.ActiveClient
var AllActiveClientsMutex	*sync.Mutex

func InitAllActiveClients() {
	AllActiveClientsMutex = &sync.Mutex{}
}

func AddAllActiveClient(addr *net.TCPAddr, id int) bool {
	AllActiveClientsMutex.Lock()
	defer AllActiveClientsMutex.Unlock()
	for _, val := range AllActiveClients {
		if val.ClientID == id {
			return false
		}
	}

	AllActiveClients = append(AllActiveClients, declarations.ActiveClient{id, addr})
	error_catcher.PushMessage(fmt.Sprintf("Added new Active User.. %s", AllActiveClients))
	return true
}

func RemoveAllActiveClientById(id int) {
	AllActiveClientsMutex.Lock()
	defer AllActiveClientsMutex.Unlock()
	for i, v := range ActiveClientsTable {
		if v.ClientID == id {
			ActiveClientsTable = append(ActiveClientsTable[:i], ActiveClientsTable[i+1:]...)
			error_catcher.PushMessage(fmt.Sprintf("Removed some user.. %s", AllActiveClients))
			return
		}
	}
}
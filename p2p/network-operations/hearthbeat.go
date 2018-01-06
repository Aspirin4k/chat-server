package network_operations

import (
	"time"
	"net"
	"fmt"
	"sync"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

var syncAddress 			*net.TCPAddr
var hearthbeatAddressMutex 	*sync.Mutex
var Changes					[]declarations.StatusChanges

func InitHearthbeat() {
	hearthbeatAddressMutex = &sync.Mutex{}
	syncAddress = nil
	go Heathbeat()
}

func SetHearthbeatAddress(addr *net.TCPAddr) {
	hearthbeatAddressMutex.Lock()
	syncAddress = addr
	hearthbeatAddressMutex.Unlock()
}

func Heathbeat() {
	for ;; {
		time.Sleep(time.Millisecond * declarations.HEARTHBEAT_TICK)
		if (syncAddress != nil) && (len((Changes)) > 0) {
			hearthbeatSend()
		}
	}
}

func hearthbeatSend() {
	hearthbeatAddressMutex.Lock()
	msg := fmt.Sprintf("%d", declarations.HEARTHBEAT)
	for _, val := range Changes {
		msg += fmt.Sprintf("\n%d %d %s", val.Status, val.ClientID, val.Address.IP.String())
	}
	SendMessage(syncAddress, msg)
	ClearChanges()
	hearthbeatAddressMutex.Unlock()
}

func AddUserOnline(addr *net.TCPAddr, id int) {
	for _, val := range Changes {
		if val.ClientID == id {
			return
		}
	}

	Changes = append(
		Changes, declarations.StatusChanges{id, addr,declarations.CLIENT_ONLINE})
}

func ClearChanges() {
	Changes = Changes[:0]
}
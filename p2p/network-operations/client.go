package network_operations

import (
	"net"
	"fmt"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

/**
Запрос для подключения клиента к одному из узлов сети
 */
func Loggining(addr *net.TCPAddr, id int, clientAddress string) {
	SendMessage(addr, fmt.Sprintf("%d %d %s", declarations.CLIENT_LOGIN, id, clientAddress))
}

func AddToOnline(addr *net.TCPAddr, id int, clientAddress string) {
	SendMessage(
		addr, fmt.Sprintf("%d %d %s", declarations.CLIENT_ADD_TO_ONLINE_CLIENTS, id, clientAddress))
}
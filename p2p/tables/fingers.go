package tables

import (
	"fmt"
	"os"
	"net"

	"github.com/sadlil/go-trigger"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

var FingerTable []declarations.Finger

/**
Добавляет новый палец в пальцевую таблицу (при условии его отсутствия)
 */
func AddFinger(node int, address *net.TCPAddr) {
	for _, v := range FingerTable {
		if v.Node == node {
			return
		}
	}

	fmt.Fprint(os.Stdout,"Adding new finger...\n")
	FingerTable = append(FingerTable, declarations.Finger{node,address})

	trigger.Fire(FINGERS_CHANGED, FingerTable)
}

/**
Находит узел или близжайшего предшественника
 */
func FindFinger(id int, serverID int) (*net.TCPAddr, bool) {
	if (id > serverID) && (id < Successor().Node) ||
		((id < Successor().Node) || (id > serverID)) && (serverID > Successor().Node) ||
		(serverID == Successor().Node) {
		return Successor().Address, true
	} else {
		return Predecessor(id), false
	}
}

/**
Получает преемника
 */
func Successor() declarations.Finger {
	return FingerTable[0]
}

/**
Получает близжайшего предшественника. Т.е. узел, идентификатор которого меньше хеша
ресурса, но больше идентификатора всех остальных таких узлов в пальцевой таблице
 */
func Predecessor(id int) *net.TCPAddr {
	num := -1
	for i, v := range FingerTable {
		if (num == -1) || (v.Node <= id) && (v.Node > FingerTable[num].Node) {
			num = i
		}
	}

	if num == -1 {
		return nil
	} else {
		return FingerTable[num].Address
	}
}
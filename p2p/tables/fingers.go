package tables

import (
	"fmt"
	"net"
	"math"
	"sort"

	"github.com/sadlil/go-trigger"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/error-catcher"
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

	error_catcher.PushMessage(
		fmt.Sprintf("Adding new finger... {%d %s}", node, address.IP.String()))
	FingerTable = append(FingerTable, declarations.Finger{node,address})

	trigger.Fire(declarations.FINGERS_CHANGED, FingerTable)
}

/**
Очищает пальцевую таблицу
 */
func ClearFingers() {
	error_catcher.PushMessage("Clearing finger table...")
	FingerTable = FingerTable[:0]
	trigger.Fire(declarations.FINGERS_CHANGED, FingerTable)
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

/**
Генерирует пальцевую таблицу
@fingers - временная таблица, необязательно размера logn
@serverID - идентификатор локального хоста
 */
func BuildFingers(fingers []declarations.Finger, serverID int) {
	ClearFingers()
	sort.Sort(ByID(fingers))
	error_catcher.PushMessage(fmt.Sprintf("Sorted temp finger table: %s",fingers))
	// Таблица должна быть log от размера хеша
	for i:=0; i<declarations.FINGERS_SIZE; i++ {
		for _, val := range fingers {
			if (len((fingers)) == 1) || ((val.Node != serverID) &&
				(serverID + int(math.Pow(2, float64(i)))) % declarations.HASH_SIZE < val.Node) {
				AddFinger(val.Node, val.Address)
				break
			}
		}
	}
}


// Кастомная сортировка
type ByID []declarations.Finger

func (s ByID) Len() int {
	return len(s)
}

func (s ByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByID) Less(i, j int) bool {
	return s[i].Node < s[j].Node
}
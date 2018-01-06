package network_operations

import (
	"net"
	"fmt"
	"strconv"

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
	message += fmt.Sprintf("\n%s", SplitResources(remoteID, serverID))
	SendMessage(address, message)
}


/**
Включает ресурсы в свои таблицы
 */
func ConcatenateResources(tokens [][]string) {
	var i int
	i = 0
	// Онлайн
	for ;(i<len(tokens)) && (tokens[i][0] != declarations.DELIM);i++ {
		val := tokens[i]
		id, _ := strconv.Atoi(val[0])
		tables.AddActiveClient(
			ParseAddress(val[1], declarations.PORT_CLIENTS), id)
	}
	i++
	// Зарегистрированные
	for ;(i<len(tokens)) && (tokens[i][0] != declarations.DELIM);i++ {
		val := tokens[i]
		id, _ := strconv.Atoi(val[0])
		key, _ := strconv.Atoi(val[2])
		tables.AddRegisteredClient(id, val[1], key)
	}
}

/**
Подготавливает ресурсы для отправки
 */
func PrepareResources() string {
	message := ""
	// Передаем подключенных к нам клиентов
	for _, val := range tables.ActiveClientsTable {
		message += fmt.Sprintf("\n%d %s", val.ClientID, val.Address.IP.String())
	}
	message += fmt.Sprintf("\n%s", declarations.DELIM)
	// Передаем зарегистрированных клиентов
	for _, val := range tables.RegisteredClientsTable {
		message += fmt.Sprintf("\n%d %s %d", val.ClientID, val.Nickname, val.Key)
	}

	if len(message) > 0 {
		return message[1:]
	} else {
		return message
	}
}

/**
Разделяет ресурсы с заданным хостом
 */
func SplitResources(remoteID int, serverID int) string {
	fmt.Println(remoteID, serverID)
	message := ""
	for i:=0; i< len(tables.ActiveClientsTable); i++ {
		v := tables.ActiveClientsTable[i]
		if (v.ClientID < remoteID) && (v.ClientID > serverID) && (remoteID > serverID) ||
			(remoteID < serverID) && ((v.ClientID > serverID) || (v.ClientID < remoteID)) {
			message += fmt.Sprintf("\n%d %s", v.ClientID, v.Address.IP.String())
			tables.RemoveActiveClientByIndex(i)
			i--
		}
	}
	message += fmt.Sprintf("\n%s", declarations.DELIM)
	for i:=0; i< len(tables.RegisteredClientsTable); i++ {
		v := tables.RegisteredClientsTable[i]
		if (v.ClientID < remoteID) && (v.ClientID > serverID) && (remoteID > serverID) ||
			(remoteID < serverID) && ((v.ClientID > serverID) || (v.ClientID < remoteID)) {
				fmt.Println(v)
			message += fmt.Sprintf("\n%d %s %d", v.ClientID, v.Nickname, v.Key)
			tables.RemoveRegisteredClientByIndex(i)
			i--
		}
	}

	if len(message) > 0 {
		return message[1:]
	} else {
		return message
	}
}
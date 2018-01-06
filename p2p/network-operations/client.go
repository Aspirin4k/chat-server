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

/**
Указание подключить данного клиента к себе
 */
func AddToOnline(addr *net.TCPAddr, id int, clientAddress string) {
	SendMessage(
		addr, fmt.Sprintf("%d %d %s", declarations.CLIENT_ADD_TO_ONLINE_CLIENTS, id, clientAddress))
}

/**
Запрос о внесении данных нового пользователя
 */
func Registering(addr *net.TCPAddr, id int, nickname string, key string, clientAddress string) {
	SendMessage(
		addr, fmt.Sprintf("%d %d %s %s %s", declarations.CLIENT_NEW, id, nickname, key, clientAddress))
}

/**
Указание добавить указанного клиента в свою таблицу зарегистрированных клиентов
 */
func AddToRegistered(addr *net.TCPAddr, id int, nickname string, key string, clientAddress string) {
	SendMessage(
		addr, fmt.Sprintf("%d %d %s %s %s", declarations.CLIENT_ADD_TO_REGISTERED_CLIENTS, id, nickname, key, clientAddress))
}
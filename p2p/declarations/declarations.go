package declarations

import "net"

// Порт, по которому будут прослушиваться другие узлы
const PORT = 7777
// Порт, по которому будут прослушиваться клиенты
const PORT_CLIENTS = 7778
const FINGERS_SIZE = 8
const HASH_SIZE = 256
const REGISTERED_CLIENTS_LOCATION = "./clients.txt"
const DELIM = "|^|"

type Finger struct {
	Node 	int
	Address *net.TCPAddr
}

type ActiveClient struct {
	ClientID int
	Address  *net.TCPAddr
}

type RegisteredClient struct {
	ClientID int
	Nickname string
	Key      int
}

/**
Команды, используемые при коммуникации между серверными узлами
 */
type Command int
const (
	NODE_JOINING Command = 1 + iota
	NODE_JOINING_ADD_BEFORE
	NODE_ADD_ME_TO_FINGER
	RESOURCE_RECEIVE_IDS
	NODE_LEAVING
	NODE_LEAVED
	FINGERS_UPDATE
	CLIENT_LOGIN
	CLIENT_ADD_TO_ONLINE_CLIENTS
	CLIENT_NEW
	CLIENT_ADD_TO_REGISTERED_CLIENTS
	UNKNOWN
)

func GetCommandByValue(command int) Command {
	switch command {
	case 1:
		return NODE_JOINING
	case 2:
		return NODE_JOINING_ADD_BEFORE
	case 3:
		return NODE_ADD_ME_TO_FINGER
	case 4:
		return RESOURCE_RECEIVE_IDS
	case 5:
		return NODE_LEAVING
	case 6:
		return NODE_LEAVED
	case 7:
		return FINGERS_UPDATE
	case 8:
		return CLIENT_LOGIN
	case 9:
		return CLIENT_ADD_TO_ONLINE_CLIENTS
	case 10:
		return CLIENT_NEW
	case 11:
		return CLIENT_ADD_TO_REGISTERED_CLIENTS
	}

	return UNKNOWN
}

/**
Команды со стороны клиента
 */
const (
	CLIENT_HELLO 	= 201
	CLIENT_REGISTER = 202
)

/**
Команды событий изменения модели
 */
const (
	FINGERS_CHANGED = "FINGERS_CHANGED"
	ACTIVE_CLIENTS_CHANGED = "ACTIVE_CLIENTS_CHANGED"
	REGISTERED_CLIENTS_CHANGED = "REGISTERED_CLIENTS_CHANGED"
)
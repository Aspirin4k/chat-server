package declarations

import "net"

const PORT = 7777

type Finger struct {
	Node 	int
	Address *net.TCPAddr
}

type Command int
const (
	NODE_JOINING Command = 1 + iota
	NODE_JOINING_ADD_BEFORE
	RESOURCE_RECEIVE_IDS
	UNKNOWN
)

func GetCommandByValue(command int) Command {
	switch command {
	case 1:
		return NODE_JOINING
	case 2:
		return NODE_JOINING_ADD_BEFORE
	case 3:
		return RESOURCE_RECEIVE_IDS
	}

	return UNKNOWN
}
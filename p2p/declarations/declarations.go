package declarations

import "net"

const PORT = 7777
const FINGERS_SIZE = 8
const HASH_SIZE = 256

type Finger struct {
	Node 	int
	Address *net.TCPAddr
}

type Command int
const (
	NODE_JOINING Command = 1 + iota
	NODE_JOINING_ADD_BEFORE
	NODE_ADD_ME_TO_FINGER
	RESOURCE_RECEIVE_IDS
	NODE_LEAVING
	NODE_LEAVED
	FINGERS_UPDATE
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
	}

	return UNKNOWN
}
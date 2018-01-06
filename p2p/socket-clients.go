package p2p

import (
	"net"
	"fmt"
	"strings"
	"strconv"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/utils"
	"github.com/Aspirin4k/chat-server/p2p/tables"
	"github.com/Aspirin4k/chat-server/p2p/network-operations"
)

func HandleClient(conn net.Conn) {
	defer conn.Close()

	error_catcher.PushMessage("New connection...")
	buffer := make([]byte, 1024)

	reqLen, err := conn.Read(buffer)
	error_catcher.CheckError(err)

	// Разбитие на токены
	input := string(buffer[:reqLen])
	error_catcher.PushMessage(fmt.Sprintf("Input data: %s", input))
	tokenRows := strings.Split(input, "\n")
	var tokens [][]string
	tokens = make([][]string, len(tokenRows), len(tokenRows))
	for i, row := range tokenRows {
		tokens[i] = strings.Split(row, " ")
	}

	// Первым идет всегда идентификатор операции
	var commandInt int
	commandInt, err = strconv.Atoi(tokens[0][0])
	error_catcher.CheckError(err)

	switch commandInt {
	// Клиент подключается к нашему узлу
	case declarations.CLIENT_HELLO:
		error_catcher.PushMessage("Client says hello!")

		hashId := utils.GetHash(tokens[0][1])
		node, isSuccessor := tables.FindFinger(hashId, ServerID)

		if isSuccessor {
			network_operations.AddToOnline(node, hashId, strings.Split(conn.RemoteAddr().String(), ":")[0])
		} else {
			network_operations.Loggining(node, hashId, strings.Split(conn.RemoteAddr().String(), ":")[0])
		}
		break
	}
}
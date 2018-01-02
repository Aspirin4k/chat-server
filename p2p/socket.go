package p2p

import (
	"net"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/network-operations"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	// Чтение всего входящего потока
	inputBytes, err := ioutil.ReadAll(conn)
	error_catcher.CheckError(err)

	// Разбитие на токены
	input := string(inputBytes)
	fmt.Fprintf(os.Stdout,"Input data: %s\n", input)
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
	command := declarations.GetCommandByValue(commandInt)

	switch command {
	// Новый узел подключается к существующей сети
	case declarations.NODE_JOINING:
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)
		// Адрес
		var remoteIP string
		remoteIP = tokens[0][2]

		fmt.Fprintf(os.Stdout,
			"New node %d %s wants to join the network\n", remoteID, remoteIP)

		nextNode, shouldAdd := tables.FindFinger(remoteID, ServerID)
		if shouldAdd {
			fmt.Fprint(os.Stdout, "Looks like he should be our successor......\n")
			network_operations.JoinAddBefore(nextNode, remoteID, remoteIP)
		} else {
			fmt.Fprint(os.Stdout, "Ask another node to take care of him...\n")
			fmt.Fprintf(os.Stdout, "Send request to %s \n", nextNode.IP.String())
			network_operations.SendMessage(nextNode, input)
		}
		break
	// Нам говорят, что перед нами должен быть добавлен новый узел
	case declarations.NODE_JOINING_ADD_BEFORE:
		// Переслать ему некоторые ресурсы
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)
		// Адрес
		var remoteIP string
		remoteIP = tokens[0][2]

		fmt.Fprintf(os.Stdout,
			"Should send to %d %s some resources before adding him to network\n", remoteID, remoteIP)

		network_operations.ReceiveIDs(network_operations.ParseAddress(tokens[0][2]), remoteID, ServerID)
		break
	case declarations.RESOURCE_RECEIVE_IDS:
		fmt.Fprint(os.Stdout, "Should add some resources ids...\n")
	}
}

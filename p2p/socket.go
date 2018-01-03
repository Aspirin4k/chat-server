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
	// Нам прислали некоторые ресурсы, с помощью которых мы должны проинициализировать таблицу ресурсов
	case declarations.RESOURCE_RECEIVE_IDS:
		fmt.Fprint(os.Stdout, "Should add some resources ids...\n")
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		// Добавим товарища в нашу пальцевую таблицу
		tables.ClearFingers()
		tables.AddFinger(remoteID,
			network_operations.ParseAddress(strings.Split(conn.RemoteAddr().String(), ":")[0]))

		// Добавление новых ресурсов
		for _, val := range tokens[1:] {
			id,_ := strconv.Atoi(val[0])
			tables.AddResource(id, network_operations.ParseAddress(val[1]))
		}

		// Просим всех обновить свои пальцевые таблицы
		network_operations.AddMeToFinger(tables.Successor().Address, ServerID, ServerAddress.IP.String())
		break
	// Нас просят добавить хост в свою пальцевую таблицу
	case declarations.NODE_ADD_ME_TO_FINGER:
		fmt.Fprint(os.Stdout, "New node in here, need to add him to our finger table...\n")
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)
		// Адрес
		var remoteIP string
		remoteIP = tokens[0][2]
		// Если к нам вернулся наш же пакет
		if remoteID == ServerID {
			fmt.Fprint(os.Stdout, "Hey! Its our packet returned to us!\n")

			var temp []declarations.Finger
			for _, val := range tokens[1:] {
				id,_ := strconv.Atoi(val[0])
				temp = append(
					temp, declarations.Finger{id, network_operations.ParseAddress(val[1])})
			}

			tables.BuildFingers(temp, ServerID)
		} else {
			fmt.Fprint(os.Stdout, "Some node somewhere joined the circle. We'll check, should we add him to our finger table!\n")

			// Добавляем новый узел в пальцевую таблицу
			var temp []declarations.Finger
			copy(tables.FingerTable, temp)
			temp = append(
				temp, declarations.Finger{remoteID, network_operations.ParseAddress(remoteIP)})

			// Оптимизируем пальцевую таблицу
			tables.BuildFingers(temp, ServerID)
			// Пересылаем сообщение дальше по кольцу
			network_operations.AddMeToFingerMessage(tables.Successor().Address, ServerID, ServerAddress.IP.String(), input)
		}
		break
	}
}

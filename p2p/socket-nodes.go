package p2p

import (
	"net"
	"fmt"
	"strings"
	"strconv"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/network-operations"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
)

func HandleConnection(conn net.Conn) {
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

		error_catcher.PushMessage(
			fmt.Sprintf("New node %d %s wants to join the network", remoteID, remoteIP))

		nextNode, shouldAdd := tables.FindFinger(remoteID, ServerID)
		if shouldAdd {
			error_catcher.PushMessage("Looks like he should be our successor...")
			error_catcher.PushMessage(
				fmt.Sprintf("Send request to %s", nextNode.IP.String()))
			network_operations.JoinAddBefore(nextNode, remoteID, remoteIP)
		} else {
			error_catcher.PushMessage("Ask another node to take care of him...")
			error_catcher.PushMessage(
				fmt.Sprintf("Send request to %s", nextNode.IP.String()))
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

		error_catcher.PushMessage(
			fmt.Sprintf("Should send to %d %s some resources before adding him to network", remoteID, remoteIP))

		network_operations.ReceiveIDs(network_operations.ParseAddress(tokens[0][2], declarations.PORT), remoteID, ServerID)
		break
	// Нам прислали некоторые ресурсы, с помощью которых мы должны проинициализировать таблицу ресурсов
	case declarations.RESOURCE_RECEIVE_IDS:
		error_catcher.PushMessage("Should add some resources ids...")
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		// Добавим товарища в нашу пальцевую таблицу
		tables.ClearFingers()
		tables.AddFinger(remoteID,
			network_operations.ParseAddress(strings.Split(conn.RemoteAddr().String(), ":")[0], declarations.PORT))

		// Добавляем его записи к своим
		network_operations.ConcatenateResources(tokens[1:])

		// Просим всех обновить свои пальцевые таблицы
		network_operations.AddMeToFinger(tables.Successor().Address, ServerID, ServerAddress.IP.String())
		break
	// Нас просят добавить хост в свою пальцевую таблицу
	case declarations.NODE_ADD_ME_TO_FINGER:
		error_catcher.PushMessage("New node in here, need to add him to our finger table...")
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)
		// Адрес
		var remoteIP string
		remoteIP = tokens[0][2]
		// Если к нам вернулся наш же пакет
		if remoteID == ServerID {
			error_catcher.PushMessage("Hey! Its our packet returned to us!")

			var temp []declarations.Finger
			for _, val := range tokens[1:] {
				id,_ := strconv.Atoi(val[0])
				temp = append(
					temp, declarations.Finger{id, network_operations.ParseAddress(val[1], declarations.PORT)})
			}

			tables.BuildFingers(temp, ServerID)
			updateHearthbeat()
		} else {
			error_catcher.PushMessage(
				"Some node somewhere joined the circle. We'll check, should we add him to our finger table!")

			// Добавляем новый узел в пальцевую таблицу
			temp := make([]declarations.Finger, len(tables.FingerTable))
			copy(temp, tables.FingerTable)
			temp = append(
				temp, declarations.Finger{remoteID, network_operations.ParseAddress(remoteIP, declarations.PORT)})

			// Оптимизируем пальцевую таблицу
			tables.BuildFingers(temp, ServerID)
			updateHearthbeat()
			// Пересылаем сообщение дальше по кольцу
			network_operations.AddMeToFingerMessage(tables.Successor().Address, ServerID, ServerAddress.IP.String(), input)
		}
		break
	// Наш предшественник покидает сеть
	case declarations.NODE_LEAVING:
		error_catcher.PushMessage("Some node leaving the circle and sending to us his resources id")

		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		// Добавляем его записи к своим
		network_operations.ConcatenateResources(tokens[1:])

		// Удаляем записи со ссылкой на него
		tables.ResourceRemoveByKey(remoteID)

		// Оповещаем остальные узлы об отключившемся
		if remoteID != tables.Successor().Node {
			network_operations.Leaved(tables.Successor().Address, ServerAddress, ServerID, remoteID)
		} else {
			// Происходит в случае, если сеть состоит из двух узлов
			tables.BuildFingers(
				[]declarations.Finger{declarations.Finger{ServerID, ServerAddress}}, ServerID)
			network_operations.SetHearthbeatAddress(nil)
		}
		break
	// Один из узлов сети отключился
	case declarations.NODE_LEAVED:
		error_catcher.PushMessage("Some node leaving the circle and sending to us his resources id")

		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		// Если отвалился наш преемник
		if tables.Successor().Node == remoteID {
			error_catcher.PushMessage("Looks like our successor leaved!")

			var temp []declarations.Finger
			for _, val := range tokens[1:] {
				id,_ := strconv.Atoi(val[0])
				temp = append(
					temp, declarations.Finger{id,
					network_operations.ParseAddress(val[1], declarations.PORT)})
			}

			tables.BuildFingers(temp, ServerID)
			updateHearthbeat()

			network_operations.UpdateFingers(
				tables.Successor().Address, ServerID, ServerAddress.IP.String(), tokens[1:])

		} else {
			// Перекличка. Добавим наш адрес в данный список
			input += fmt.Sprintf("\n%d %s", ServerID, ServerAddress.IP.String())
			network_operations.SendMessage(tables.Successor().Address, input)
		}

		tables.ResourceRemoveByKey(remoteID)
		break
	// Обновление пальцевых таблиц
	case declarations.FINGERS_UPDATE:
		error_catcher.PushMessage("Stabilization after node leaved!")

		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		// Если это не наш пакет, то обновляем пальцевую таблицу и пересылаем его далее
		if remoteID != ServerID {

			var temp []declarations.Finger
			for _, val := range tokens[1:] {
				id,_ := strconv.Atoi(val[0])
				temp = append(
					temp,
					declarations.Finger{id, network_operations.ParseAddress(val[1], declarations.PORT)})
			}

			tables.BuildFingers(temp, ServerID)
			updateHearthbeat()
			network_operations.SendMessage(tables.Successor().Address, input)
		}
		break
	// Подключается клиент, необходимо найти хост, который будет его слушать
	case declarations.CLIENT_LOGIN:
		error_catcher.PushMessage("Someone asks us to join client..")

		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		node, isSuccessor := tables.FindFinger(remoteID, ServerID)

		if isSuccessor {
			network_operations.AddToOnline(node, remoteID, tokens[0][2])
		} else {
			network_operations.SendMessage(node, input)
		}
		break
	// К нашему хосту должен подключиться клиент
	case declarations.CLIENT_ADD_TO_ONLINE_CLIENTS:
		error_catcher.PushMessage("We should join client..")
		// TODO: ПРОВЕРИТЬ КЛИЕНТА НА СУЩЕСТВОВАНИЕ
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		tables.AddActiveClient(
			network_operations.ParseAddress(tokens[0][2], declarations.PORT_CLIENTS), remoteID)
		network_operations.AddUserOnline(
			network_operations.ParseAddress(tokens[0][2], declarations.PORT_CLIENTS), remoteID)

		// TODO: ОТВЕТИТЬ КЛИЕНТУ
		break
	// Попросить узел найти того, кто может добавить пользователя в таблицу
	// зарегистрированных пользователей
	case declarations.CLIENT_NEW:
		error_catcher.PushMessage("Someone ask no register new client")
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		node, isSuccessor := tables.FindFinger(remoteID, ServerID)

		if isSuccessor {
			network_operations.AddToRegistered(
				node, remoteID, tokens[0][2], tokens[0][3], tokens[0][4])
		} else {
			network_operations.SendMessage(node, input)
		}
		break
	// Добавить пользователя в нашу таблицу зарегистрированных пользователей
	case declarations.CLIENT_ADD_TO_REGISTERED_CLIENTS:
		error_catcher.PushMessage("We should register client..")
		// TODO: ПРОВЕРИТЬ КЛИЕНТА НА СУЩЕСТВОВАНИЕ
		// Идентификатор удаленного узла
		var remoteID int
		remoteID, err = strconv.Atoi(tokens[0][1])
		error_catcher.CheckError(err)

		var key int
		key, err = strconv.Atoi(tokens[0][3])
		error_catcher.CheckError(err)

		tables.AddRegisteredClient(remoteID, tokens[0][2], key)
		break
	// Синхронизация общих параметров
	case declarations.HEARTHBEAT:
		error_catcher.PushMessage("Synchronization..")
		for _, val := range tokens[1:] {
			status, _ := strconv.Atoi(val[0])
			id, _ := strconv.Atoi(val[1])

			switch status {
			case declarations.CLIENT_ONLINE:
				if tables.AddAllActiveClient(
						network_operations.ParseAddress(val[2], declarations.PORT_CLIENTS), id) {
					network_operations.AddUserOnline(
						network_operations.ParseAddress(val[2], declarations.PORT_CLIENTS), id)
				}
				break
			case declarations.CLIENT_OFFLINE:
				tables.RemoveAllActiveClientById(id)
				break
			}
		}
		break
	}
}

func updateHearthbeat() {
	// Изменяем адрес, который му будем оповещать о событиях
	if tables.Successor().Node != ServerID {
		network_operations.SetHearthbeatAddress(tables.Successor().Address)
		for _, val := range tables.AllActiveClients {
			network_operations.AddUserOnline(val.Address, val.ClientID)
		}
	} else {
		network_operations.SetHearthbeatAddress(nil)
	}
}
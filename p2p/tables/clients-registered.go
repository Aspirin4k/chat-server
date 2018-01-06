package tables

import (
	"os"
	"bufio"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/error-catcher"

	"github.com/sadlil/go-trigger"
	"strconv"
)

var RegisteredClientsTable []declarations.RegisteredClient

func InitRegisteredClients() {
	RegisteredClientsTable = LoadRegisteredClients([]declarations.RegisteredClient{})

	trigger.Fire(declarations.REGISTERED_CLIENTS_CHANGED, RegisteredClientsTable)
}

func AddRegisteredClient(id int, nickname string, key int) bool {
	for _, val := range RegisteredClientsTable {
		if val.ClientID == id {
			return false
		}
	}

	RegisteredClientsTable =
		append(RegisteredClientsTable, declarations.RegisteredClient{id,nickname, key})

	trigger.Fire(declarations.REGISTERED_CLIENTS_CHANGED, RegisteredClientsTable)
	return true
}

func RemoveRegisteredClientByKey(id int) {
	for i, v := range RegisteredClientsTable {
		if v.ClientID == id {
			RegisteredClientsTable = append(RegisteredClientsTable[:i], RegisteredClientsTable[i+1:]...)
		}
	}

	trigger.Fire(declarations.REGISTERED_CLIENTS_CHANGED, RegisteredClientsTable)
}

func RemoveRegisteredClientByIndex(i int) {
	RegisteredClientsTable = append(RegisteredClientsTable[:i], RegisteredClientsTable[i+1:]...)

	trigger.Fire(declarations.REGISTERED_CLIENTS_CHANGED, RegisteredClientsTable)
}

func SaveRegisteredClients() {
	if len(RegisteredClientsTable) <= 0 {
		return
	}
	// Чтобы не перезаписать уже записанных пользователей
	buffer := LoadRegisteredClients(RegisteredClientsTable)

	f, err := os.Create(declarations.REGISTERED_CLIENTS_LOCATION)
	error_catcher.CheckError(err)

	defer f.Close()

	w := bufio.NewWriter(f)
	for _, val := range buffer {
		w.WriteString(fmt.Sprintf("%d %s %d\n", val.ClientID, val.Nickname, val.Key))
	}
	w.Flush()
}

/**
Подгружает новых клиентов из файла
 */
func LoadRegisteredClients(current []declarations.RegisteredClient) []declarations.RegisteredClient {
	buffer := make([]declarations.RegisteredClient, len(current))
	copy(buffer, current)

	if _, err := os.Stat(declarations.REGISTERED_CLIENTS_LOCATION); err == nil {
		b, err := ioutil.ReadFile(declarations.REGISTERED_CLIENTS_LOCATION)
		error_catcher.CheckError(err)
		lines := strings.Split(string(b),"\n")
		for _, val := range lines {
			tokens := strings.Split(val, " ")

			if len(tokens) == 3 {
				// Проверка на существование такой записи у нас (незачем повторяться)
				id, _ := strconv.Atoi(tokens[0])
				var exist bool
				exist = false
				for _, bufferVal := range buffer {
					if bufferVal.ClientID == id {
						exist = true
						break
					}
				}

				if !exist {
					key, _ := strconv.Atoi(tokens[2])

					buffer = append(
						buffer, declarations.RegisteredClient{id, tokens[1], key})
				}
			}
		}
	}

	return buffer
}
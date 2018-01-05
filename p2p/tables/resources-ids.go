package tables

import (
	"net"
	"fmt"
	"os"
)

type ResourceID struct {
	ID 		int
	HostID  int
	Address *net.TCPAddr
}
var ResourcesIDsTable []ResourceID

/**
Добавляет ссылку на ресурс в таблицу идентификаторов ресурсов
 */
func AddResource(id int, ownerId int, owner *net.TCPAddr) {
	for _, v := range ResourcesIDsTable {
		if v.ID == id {
			return
		}
	}

	fmt.Fprint(os.Stdout,"Adding new resource identificator...\n")
	ResourcesIDsTable = append(ResourcesIDsTable, ResourceID{id, ownerId, owner})
}

/**
Получает адрес узла, на котором расположен данный ресурс
 */
func ResourceAddressByKey(id int) *net.TCPAddr {
	for _, v := range ResourcesIDsTable {
		if v.ID == id {
			return v.Address
		}
	}
	return nil
}

/**
Удаляет ссылку на ресурс
 */
func ResourceRemoveByKey(id int) {
	for i, v := range ResourcesIDsTable {
		if v.ID == id {
			ResourcesIDsTable = append(ResourcesIDsTable[:i], ResourcesIDsTable[i+1:]...)
		}
	}
}
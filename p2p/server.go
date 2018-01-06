package p2p

import (
	"net"
	"fmt"
	"os"
	"os/signal"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/tables"
	"github.com/Aspirin4k/chat-server/p2p/network-operations"
	"github.com/Aspirin4k/chat-server/p2p/utils"
)

var ServerID 	  			int
var ServerAddress 			*net.TCPAddr
var ClientListenerAddress 	*net.TCPAddr
var NodesPort               int
var ClientsPort             int

func CreateAndListen(addr net.IP, nodesPort int, clientsPort int, remoteIP net.IP) {
	// Генерируем идентификатор и адрес
	Create(addr, nodesPort, clientsPort)

	// Отправляем запрос на подключение к сети
	if remoteIP != nil {
		remoteAddress, err := net.ResolveTCPAddr("tcp4",
			fmt.Sprintf("%s:%d", remoteIP.String(), NodesPort))
		error_catcher.CheckError(err)
		network_operations.Join(remoteAddress, ServerID, ServerAddress)
	}

	go listenNodes()
	listenClients()
}

func Create(addr net.IP, nodesPort int, clientsPort int) {
	var err error
	NodesPort = nodesPort
	ServerAddress, err = net.ResolveTCPAddr("tcp4",
		fmt.Sprintf("%s:%d", addr.String(), NodesPort))
	error_catcher.CheckError(err)

	ClientsPort = clientsPort
	ClientListenerAddress, err = net.ResolveTCPAddr("tcp4",
		fmt.Sprintf("%s:%d", addr.String(), ClientsPort))

	// Идентификатором узла является хеш от адреса
	ServerID = utils.GetHash(addr.String())

	tables.Init()
	// Первая запись в пальцевой таблице
	tables.AddFinger(ServerID, ServerAddress)

	// Вешаем запрос при отключении хоста
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for _ = range signalChannel {
			if len(tables.FingerTable) > 0 {
				network_operations.Leave(tables.Successor().Address, tables.Successor().Node, ServerID)
			}
			os.Exit(0)
		}
	}()
}

/**
Устанавливает прослушивание других серверных узлов
 */
func listenNodes() {
	listener, err := net.ListenTCP("tcp", ServerAddress)
	error_catcher.CheckError(err)

	error_catcher.PushMessage("Starting listening nodes...")
	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go HandleConnection(conn)
	}
}

/**
Устанавливает прослушивание клиентских приложений
 */
func listenClients() {
	listener, err := net.ListenTCP("tcp", ClientListenerAddress)
	error_catcher.CheckError(err)

	error_catcher.PushMessage("Starting listening clients...")
	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go HandleClient(conn)
	}
}
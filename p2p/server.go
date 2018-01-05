package p2p

import (
	"net"
	"fmt"
	"os"
	"crypto/sha1"
	"io"
	"os/signal"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
	"github.com/Aspirin4k/chat-server/p2p/network-operations"
)

var ServerID 	  int
var ServerAddress *net.TCPAddr

func CreateAndListen(addr net.IP, remoteIP net.IP) {
	// Генерируем идентификатор и адрес
	Create(addr)

	listener, err := net.ListenTCP("tcp", ServerAddress)
	error_catcher.CheckError(err)

	// Вешаем запрос при отключении хоста
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		for _ = range signalChannel {
			network_operations.Leave(tables.Successor().Address, tables.Successor().Node, ServerID)
			os.Exit(0)
		}
	}()

	// Отправляем запрос на подключение к сети
	if remoteIP != nil {
		remoteAddress, err := net.ResolveTCPAddr("tcp4",
							fmt.Sprintf("%s:%d", remoteIP.String(), declarations.PORT))
		error_catcher.CheckError(err)
		network_operations.Join(remoteAddress, ServerID, ServerAddress)
	}

	error_catcher.PushMessage("Starting listening...")
	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go HandleConnection(conn)
	}
}

func Create(addr net.IP) {
	var err error
	ServerAddress, err = net.ResolveTCPAddr("tcp4",
		fmt.Sprintf("%s:%d", addr.String(), declarations.PORT))
	error_catcher.CheckError(err)

	// Идентификатором узла является хеш от адреса
	sha := sha1.New()
	io.WriteString(sha, addr.String())
	ServerID = int(sha.Sum(nil)[0])

	tables.Init()
	// Первая запись в пальцевой таблице
	tables.AddFinger(ServerID, ServerAddress)
}
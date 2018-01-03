package p2p

import (
	"net"
	"fmt"
	"os"
	"crypto/sha1"
	"io"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/p2p/tables"
)

var ServerID 	  int
var ServerAddress *net.TCPAddr

func CreateAndListen(addr net.IP) {
	// Генерируем идентификатор и адрес
	Create(addr)

	listener, err := net.ListenTCP("tcp", ServerAddress)
	error_catcher.CheckError(err)

	fmt.Fprint(os.Stdout,"Starting listening...\n")
	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		fmt.Fprint(os.Stdout,"New connection...\n")
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
	// tables.AddFinger(ServerID, ServerAddress)
}
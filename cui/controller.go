package cui

import (
	"time"

	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"github.com/Aspirin4k/chat-server/error-catcher"
)

func UpdateChanges(event string) interface{} {
	for (g == nil) && (fingers == nil) {
		time.Sleep(time.Millisecond * 500)
	}

	switch event {
	case declarations.FINGERS_CHANGED:
		return func(table []declarations.Finger) {
			updateFingerTable(fingers, table, g)
		}
	case declarations.ACTIVE_CLIENTS_CHANGED:
		return func(table []declarations.ActiveClient) {
			updateActiveClientsTable(activeClients, table, g)
		}
	}

	return nil
}

func listenForMessages() {
	buffer := make([]byte, 1024)
	for {
		length, err := error_catcher.LocalReader.Read(buffer)
		error_catcher.CheckError(err)

		msg := string(buffer[:length])
		addMessage(messages, msg, g)
	}
}
package tables

import (
	"github.com/sadlil/go-trigger"

	"github.com/Aspirin4k/chat-server/cui"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

func Init() {
	trigger.On(
		declarations.FINGERS_CHANGED, cui.UpdateChanges(declarations.FINGERS_CHANGED))
	trigger.On(
		declarations.ACTIVE_CLIENTS_CHANGED, cui.UpdateChanges(declarations.ACTIVE_CLIENTS_CHANGED))
}
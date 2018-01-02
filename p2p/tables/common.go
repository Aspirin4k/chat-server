package tables

import (
	"github.com/sadlil/go-trigger"

	"github.com/Aspirin4k/chat-server/cui"
)

const (
	FINGERS_CHANGED = "FINGERS_CHANGED"
)

func Init() {
	trigger.On(FINGERS_CHANGED, cui.UpdateChanges(FINGERS_CHANGED))
}
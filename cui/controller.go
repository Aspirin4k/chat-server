package cui

import (
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

func UpdateChanges(event string) interface{} {
	return func(fingers []declarations.Finger) {
		rerender(fingers)
	}
}

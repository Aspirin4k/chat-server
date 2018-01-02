package cui

import (
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"fmt"
)

func UpdateChanges(event string) interface{} {
	return func(fingers []declarations.Finger) {
		fmt.Println(fingers)
		rerender(fingers)
	}
}

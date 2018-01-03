package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"time"
	"os"
	"sync"
)

var g *gocui.Gui

var (
	fingersMutex *sync.Mutex
)

/**
Функция входа. Устанавливает нужные лейауты.
 */
func Render() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	error_catcher.CheckError(err)
	defer g.Close()

	fingersMutex = &sync.Mutex{}

	g.SetManagerFunc(layout)

	err = keybindings(g)
	error_catcher.CheckError(err)

	err = g.MainLoop()
	error_catcher.CheckError(err)
}

/**
Основной лейаут
 */
func layout(g *gocui.Gui) error {
	if _, err := g.SetView("fingers", 0, 0, 22, 8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func rerender(val []declarations.Finger) {
	for g == nil {
		time.Sleep(time.Millisecond * 500)
	}

	fingersMutex.Lock()
	g.Update(func(g *gocui.Gui) error {
		view, err := g.View("fingers")
		if err != nil {
			fingersMutex.Unlock()
			return err
		}

		view.Clear()
		fmt.Fprintf(os.Stdout, "Rendering %s\n", val)
		for _, finger := range val {
			fmt.Fprintf(view, "%4d %15s", finger.Node, finger.Address.IP.String())
		}
		fingersMutex.Unlock()
		return nil
	})
}
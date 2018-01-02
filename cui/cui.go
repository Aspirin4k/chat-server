package cui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
	"time"
)

var g *gocui.Gui

/**
Функция входа. Устанавливает нужные лейауты.
 */
func Render() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	error_catcher.CheckError(err)
	defer g.Close()

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
	if _, err := g.SetView("fingers", 0, 0, 30, 8); err != nil {
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

	g.Update(func(g *gocui.Gui) error {
		view, err := g.View("fingers")
		if err != nil {
			return err
		}

		view.Clear()
		for _, finger := range val {
			fmt.Fprintf(view, "%d %s", finger.Node, finger.Address.IP.String())
		}
		return nil
	})
}
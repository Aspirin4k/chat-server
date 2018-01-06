package cui

import (
	"github.com/jroimartin/gocui"
	"github.com/Aspirin4k/chat-server/error-catcher"
)

var g *gocui.Gui

var (
	messages 			*MessagesWidget
	fingers 			*FingersWidget
	activeClients		*ActiveClientsWidget
	registeredClients 	*RegisteredClientsWidget
)

/**
Функция входа. Устанавливает нужные лейауты.
 */
func Render() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	error_catcher.CheckError(err)
	defer g.Close()

	messages = NewMessagesWidget(0, 9, 79, 15,"")
	fingers = NewFingersWidget(0, 0, 22, 8)
	activeClients = NewActiveClientsWidget(23,0,22,8)
	registeredClients = NewRegisteredClientsWidget(46, 0, 33, 8)

	g.SetManager(messages, fingers, activeClients, registeredClients)

	err = keybindings(g)
	error_catcher.CheckError(err)
	// Слушаем поток сообщений
	go listenForMessages()

	err = g.MainLoop()
	if err != gocui.ErrQuit {
		error_catcher.CheckError(err)
	}
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding(
		"", gocui.KeyArrowUp, gocui.ModNone, messagesMoveUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(
		"", gocui.KeyArrowDown, gocui.ModNone, messagesMoveDown); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func messagesMoveUp(g *gocui.Gui, v *gocui.View) error {
	setMessagesPos(messages, messages.pos - 1)
	return nil
}

func messagesMoveDown(g *gocui.Gui, v *gocui.View) error {
	setMessagesPos(messages, messages.pos + 1)
	return nil
}

func setMessagesPos(m *MessagesWidget, newPos int) error {
	if (newPos > 0) && (newPos + m.h <= len(m.lines) + 1) {
		m.pos = newPos
	}
	return nil
}
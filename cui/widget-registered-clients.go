package cui

import (
	"sync"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

type RegisteredClientsWidget struct {
	x, y int
	w, h int
	lines []string
	mutex *sync.Mutex
}

func NewRegisteredClientsWidget(x,y int, w,h int) *RegisteredClientsWidget {
	return &RegisteredClientsWidget{x,y,w,h, []string{}, &sync.Mutex{}}
}

func (w *RegisteredClientsWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView("registered-clients", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != gocui.ErrUnknownView {
		error_catcher.CheckError(err)
	}

	view.Title = "All"
	w.render(view)
	return nil
}

/**
Устанавливает значения строк. Вычисляет строки так, чтобы поместились в заданную ширину
 */
func (w *RegisteredClientsWidget) SetVal(lines []declarations.RegisteredClient) {
	w.lines = w.lines[:0]
	for _, client := range lines {
		w.lines = append(
			w.lines, fmt.Sprintf("%4d %18s %8d", client.ClientID, client.Nickname, client.Key))
	}
}

func (w *RegisteredClientsWidget) Val() []string {
	return w.lines
}

func (w *RegisteredClientsWidget) render(view *gocui.View) {
	view.Clear()
	for _, line := range w.lines {
		fmt.Fprintf(view, "%s\n", line)
	}
}

func (w *RegisteredClientsWidget) update(g *gocui.Gui) error {
	view, err := g.View("registered-clients")
	if err != nil {
		w.mutex.Unlock()
		return err
	}

	w.render(view)
	w.mutex.Unlock()
	return nil
}

func updateRegisteredClientsTable(client *RegisteredClientsWidget, table []declarations.RegisteredClient, gui *gocui.Gui) {
	client.mutex.Lock()
	client.SetVal(table)
	gui.Update(client.update)
}
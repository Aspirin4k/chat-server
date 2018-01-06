package cui

import (
	"sync"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/Aspirin4k/chat-server/error-catcher"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

type ActiveClientsWidget struct {
	x, y int
	w, h int
	lines []string
	mutex *sync.Mutex
}

func NewActiveClientsWidget(x,y int, w,h int) *ActiveClientsWidget {
	return &ActiveClientsWidget{x,y,w,h, []string{}, &sync.Mutex{}}
}

func (w *ActiveClientsWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView("active-clients", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != gocui.ErrUnknownView {
		error_catcher.CheckError(err)
	}

	view.Title = "Online"
	w.render(view)
	return nil
}

/**
Устанавливает значения строк. Вычисляет строки так, чтобы поместились в заданную ширину
 */
func (w *ActiveClientsWidget) SetVal(lines []declarations.ActiveClient) {
	w.lines = w.lines[:0]
	for _, client := range lines {
		w.lines = append(
			w.lines, fmt.Sprintf("%4d %15s", client.ClientID, client.Address.IP.String()))
	}
}

func (w *ActiveClientsWidget) Val() []string {
	return w.lines
}

func (w *ActiveClientsWidget) render(view *gocui.View) {
	view.Clear()
	for _, line := range w.lines {
		fmt.Fprintf(view, "\u001b[32m%s\u001b[0m\n", line)
	}
}

func (w *ActiveClientsWidget) update(g *gocui.Gui) error {
	view, err := g.View("active-clients")
	if err != nil {
		w.mutex.Unlock()
		return err
	}

	w.render(view)
	w.mutex.Unlock()
	return nil
}

func updateActiveClientsTable(client *ActiveClientsWidget, table []declarations.ActiveClient, gui *gocui.Gui) {
	client.mutex.Lock()
	client.SetVal(table)
	gui.Update(client.update)
}
package cui

import (
	"sync"
	"github.com/jroimartin/gocui"
	"github.com/Aspirin4k/chat-server/error-catcher"
	"fmt"
	"github.com/Aspirin4k/chat-server/p2p/declarations"
)

type FingersWidget struct {
	x, y int
	w, h int
	lines []string
	mutex *sync.Mutex
}

func NewFingersWidget(x,y int, w,h int) *FingersWidget {
	return &FingersWidget{x,y,w,h, []string{}, &sync.Mutex{}}
}

func (w *FingersWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView("fingers", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != gocui.ErrUnknownView {
		error_catcher.CheckError(err)
	}

	view.Title = "Fingers"
	w.render(view)
	return nil
}

/**
Устанавливает значения строк. Вычисляет строки так, чтобы поместились в заданную ширину
 */
func (w *FingersWidget) SetVal(lines []declarations.Finger) {
	w.lines = w.lines[:0]
	for _, finger := range lines {
		w.lines = append(
			w.lines, fmt.Sprintf("%4d %15s", finger.Node, finger.Address.IP.String()))
	}
}

func (w *FingersWidget) Val() []string {
	return w.lines
}

func (w *FingersWidget) render(view *gocui.View) {
	view.Clear()
	for _, line := range w.lines {
		fmt.Fprintf(view, "%s\n", line)
	}
}

func (w *FingersWidget) update(g *gocui.Gui) error {
	view, err := g.View("fingers")
	if err != nil {
		w.mutex.Unlock()
		return err
	}

	w.render(view)
	w.mutex.Unlock()
	return nil
}

func updateFingerTable(finger *FingersWidget, table []declarations.Finger, gui *gocui.Gui) {
	finger.mutex.Lock()
	finger.SetVal(table)
	gui.Update(finger.update)
}
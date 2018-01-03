package cui

import (
	"strings"
	"github.com/jroimartin/gocui"
	"github.com/Aspirin4k/chat-server/error-catcher"
	"fmt"
	"sync"
)

type MessagesWidget struct {
	x, y int
	w, h int
	lines []string
	pos int
	mutex *sync.Mutex
}

/**
Создает дефолтный виджет с заданным сообщением
@x,y - координаты
@w,h - ширина и высота. Строки обрезаются по ширине
 */
func NewMessagesWidget(x,y int, w,h int, body string) *MessagesWidget {
	messages := &MessagesWidget{x,y,w,h,[]string{},0, &sync.Mutex{}}

	lines := strings.Split(body, "\n")
	messages.SetVal(lines)

	return messages
}

/**
Лайоут. Стандартный рендер строк, которые находятся в данной зоне видимости
 */
func (w *MessagesWidget) Layout(g *gocui.Gui) error {
	view, err := g.SetView("messages", w.x, w.y, w.x+w.w, w.y+w.h)
	if err != gocui.ErrUnknownView {
		error_catcher.CheckError(err)
	}

	view.Title = "Logged information"
	w.render(view)
	return nil
}

/**
Устанавливает значения строк. Вычисляет строки так, чтобы поместились в заданную ширину
 */
func (w *MessagesWidget) SetVal(lines []string) {
	pos := 0
	for i:=0; i<len(lines); i++ {
		// Если строка не помещается в блок, то обрезаем ее на несколько
		if len(lines[i]) > w.w {
			buffer := lines[i]
			lines[i] = buffer[:w.w]
			lines = append(lines[:i+1], append([]string{buffer[w.w:]},lines[i+1:]...)...)
		}

		if len(lines[i]) == 0 {
			lines = append(lines[:i], lines[i+1:]...)
			i--
		}

		// Если не помещаемся по высоте, то изменяем позицию текущей строки
		if i > pos + w.h - 1 {
			pos = i - w.h + 1
		}
	}

	w.pos = pos
	w.lines = lines
}

func (w *MessagesWidget) Val() []string {
	return w.lines
}

func (w *MessagesWidget) render(view *gocui.View) {
	view.Clear()
	for i:=w.pos; (i<w.pos + w.h) && (i<len(w.lines)); i++ {
		fmt.Fprintf(view, "%s\n", w.lines[i])
	}
}

/**
Перерисовка виджета
 */
func (w *MessagesWidget) update(g *gocui.Gui) error {
	w.mutex.Lock()
	view, err := g.View("messages")
	if err != nil {
		w.mutex.Unlock()
		return err
	}

	w.render(view)
	w.mutex.Unlock()
	return nil
}

/**
Добовляет сообщение в виджет и перерисовывает его
 */
func addMessage(message *MessagesWidget, msg string, gui *gocui.Gui) {
	newLines := make([]string, len(message.Val()))
	copy(newLines, message.Val())
	msgSplitted := strings.Split(msg, "\n")
	for _, msgSplittedToken := range msgSplitted {
		newLines = append(newLines, msgSplittedToken)
	}
	message.SetVal(newLines)
	gui.Update(messages.update)
}
package cli

import "github.com/nsf/termbox-go"

type List struct {
	*view
	lines []string
}

func (l *List) AddLine(line string) {
	l.lines = append(l.lines, line)
}

func (l *List) Render() {
	l.Clear()

	for y, line := range l.lines {
		for i, ch := range line {
			termbox.SetCell(l.x+i, l.y+y, ch, fgColor, bgColor)
		}
	}
}

func NewList() *List {
	return &List{view: BlankView()}
}

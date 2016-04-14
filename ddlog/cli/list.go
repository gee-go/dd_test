package cli

import "github.com/nsf/termbox-go"

type Line struct {
	*view
	val string
}

func NewLine() *Line {
	return &Line{view: BlankView()}
}

func (l *Line) Set(v string) {
	if l.val == v {
		return
	}

	l.val = v
	j := 0
	for _, ch := range l.val {
		termbox.SetCell(l.x+j, l.y, ch, fgColor, bgColor)
		j++
	}
}

type List struct {
	*view
	lines []string
}

func (l *List) AddLine(line string) {
	l.lines = append(l.lines, line)
	over := len(l.lines) - l.h

	if over > 0 {
		l.lines = l.lines[over:]
	}
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

package cli

import "github.com/nsf/termbox-go"

const (
	fgColor = termbox.ColorWhite
	bgColor = termbox.ColorDefault
)

type view struct {
	x, y int
	w, h int
}

func BlankView() *view {
	return NewView(0, 0, 0, 0)
}

func NewView(x, y, w, h int) *view {
	return &view{x, y, w, h}
}

func (v *view) Clear() {
	for i := 0; i < v.w; i++ {
		for j := 0; j < v.h; j++ {
			termbox.SetCell(v.x+i, v.y+j, ' ', fgColor, bgColor)
		}
	}
}

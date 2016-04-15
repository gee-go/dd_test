package cli

import (
	"bytes"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type Row struct {
	Cells []string
}

func NewRow(cells ...string) *Row {
	return &Row{Cells: cells}
}

func (r *Row) AddCol(v string) {
	r.Cells = append(r.Cells, v)
}

type Table struct {
	*view

	maxColWidth []int
	colXOffset  []int
	padding     int

	rows []*Row
}

func (tk *Table) Width() int {
	if len(tk.rows) == 0 {
		return 0
	}

	w := (len(tk.rows[0].Cells) + 1) * tk.padding
	for _, cw := range tk.maxColWidth {
		w += cw
	}

	return w
}

func (tk *Table) ResetRows() {
	tk.rows = make([]*Row, 0)
}

func (tk *Table) AddRow(row *Row) {
	tk.rows = append(tk.rows, row)
}

// updateMaxColWidth keeps track of the widest col we have seen so far.
// Only grow columns to avoid jumpiness
func (tk *Table) updateMaxColWidth(parts []string) {

	for i, part := range parts {
		w := runewidth.StringWidth(part)

		if i > len(tk.maxColWidth)-1 {
			tk.maxColWidth = append(tk.maxColWidth, w)
		} else if w > tk.maxColWidth[i] {
			tk.maxColWidth[i] = w
		}
	}
}

// Given the current max col widths and padding, where should a column start.
func (tk *Table) updateColXOffset() {
	tk.colXOffset = make([]int, len(tk.maxColWidth))

	off := tk.padding
	for i, w := range tk.maxColWidth {
		tk.colXOffset[i] = off
		off += tk.padding + w
	}
}

func (tk *Table) prerender() {
	for _, row := range tk.rows {
		tk.updateMaxColWidth(row.Cells)
	}

	tk.updateColXOffset()
}

func (tk *Table) String() string {
	tk.prerender()

	var b bytes.Buffer
	for _, row := range tk.rows {
		current := 0
		for c, col := range row.Cells {
			x := tk.colXOffset[c]
			for current < x {
				b.WriteRune(' ')
				current++
			}
			current += runewidth.StringWidth(col)
			b.WriteString(col)
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (tk *Table) Render() {
	tk.prerender()
	tk.Clear()
	for r, row := range tk.rows {
		for c, col := range row.Cells {
			x := tk.colXOffset[c]
			for i, ch := range col {
				termbox.SetCell(tk.x+x+i, tk.y+r, ch, fgColor, bgColor)
			}
		}
	}

}

func NewTable() *Table {

	return &Table{
		padding: 5,
		view:    BlankView(),
	}
}

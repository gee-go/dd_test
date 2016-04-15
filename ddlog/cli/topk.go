package cli

import (
	"fmt"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/gee-go/ddlog/ddlog"
)

func SetTopK(t *Table, pages []*ddlog.PageCount) {
	t.ResetRows()
	head := NewRow("Hits", "Hit %", "Bytes", "Page")
	t.AddRow(head)

	for _, page := range pages {
		row := NewRow()
		row.AddCol(strconv.Itoa(page.Count))
		row.AddCol(fmt.Sprintf("%3.1f%%", 100*page.CountPercent))
		row.AddCol(humanize.Bytes(page.Bytes))
		row.AddCol(page.Name)
		t.AddRow(row)
	}
}

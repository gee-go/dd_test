package cli

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gee-go/ddlog/ddlog"
	"github.com/nsf/termbox-go"
)

type UI struct {
	Mon       *ddlog.Monitor
	TopKTable *Table
	AlertList *List

	quitChan   chan bool
	resizeChan chan termbox.Event
}

func NewUI(mon *ddlog.Monitor) *UI {
	return &UI{
		Mon: mon,

		quitChan:   make(chan bool, 1),
		resizeChan: make(chan termbox.Event),
		TopKTable:  NewTable(),
		AlertList:  NewList(),
	}
}

func (ui *UI) UpdateTopK(k int) {
	pages := ui.Mon.TopK(k)
	ui.TopKTable.ResetRows()
	head := NewRow("Hits", "Hit %", "Bytes", "Page")
	ui.TopKTable.AddRow(head)

	for _, page := range pages {
		row := NewRow()
		row.AddCol(strconv.Itoa(page.Count))
		row.AddCol(fmt.Sprintf("%3.1f%%", 100*page.CountPercent))
		row.AddCol(humanize.Bytes(page.Bytes))
		row.AddCol(page.Name)
		ui.TopKTable.AddRow(row)
	}

	ui.TopKTable.Render()
}

func (ui *UI) UpdateAlert() {
	ui.AlertList.AddLine("aaa")
	ui.AlertList.Render()
}

func (ui *UI) Resize() {
	w, h := termbox.Size()

	headerHeight := 3
	ui.TopKTable.w = w / 2
	ui.TopKTable.y = headerHeight
	ui.TopKTable.h = h - headerHeight

	ui.AlertList.w = w / 2
	ui.AlertList.y = headerHeight
	ui.AlertList.h = h - headerHeight
	ui.AlertList.x = w / 2
}

func (ui *UI) StartUpdate(rate time.Duration) {
	refreshTicker := time.NewTicker(rate)
	defer refreshTicker.Stop()

	ui.Resize()

	for {
		select {
		case <-ui.quitChan:
			return
		case <-ui.resizeChan:
			ui.Resize()
		case <-refreshTicker.C:
			ui.UpdateAlert()
			ui.UpdateTopK(ui.TopKTable.h)
			termbox.Flush()
		}
	}
}

func (ui *UI) Start() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()
	go ui.StartUpdate(1 * time.Second)
	for {
		ev := termbox.PollEvent()

		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlC, termbox.KeyEsc:
				return
			}
		case termbox.EventError:
			log.Fatal(ev.Err)
		case termbox.EventResize:
			select {
			case ui.resizeChan <- ev:
				// non blocking send
			default:
				// go on
			}

		}
	}
}

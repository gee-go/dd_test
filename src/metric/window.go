package metric

import (
	"time"

	"github.com/gee-go/dd_test/ddlog"
	"github.com/k0kubun/pp"
)

type Counter struct {
	Total  int
	ByPage map[string]int
}

func (c *Counter) Inc(page string) {
	c.Total++
	if _, found := c.ByPage[page]; !found {
		c.ByPage[page] = 0
	}

	c.ByPage[page]++
}

func NewCounter() *Counter {
	return &Counter{
		ByPage: make(map[string]int),
	}
}

type RollingCounter struct {
	data []*Counter

	current int
}

func NewRollingCounter(size int) *RollingCounter {
	return &RollingCounter{
		data: make([]*Counter, size),
	}
}

func (w *RollingCounter) Step() {
	// move current, and clear old
	w.current = (w.current + 1) % len(w.data)
	w.data[w.current] = NewCounter()
}

func (w *RollingCounter) Count(m *ddlog.Message) {
	c := w.data[w.current]
	if c == nil {
		c = NewCounter()
		w.data[w.current] = c
	}
	c.Inc(m.EventName())
}

type Group struct {
	count10s *RollingCounter
}

func (g *Group) Start(lines chan *ddlog.Message) {
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case m := <-lines:
			g.count10s.Count(m)
		case <-tick:
			g.count10s.Step()
			pp.Println(g.count10s)
		}
	}

}

func NewGroup() *Group {
	return &Group{
		count10s: NewRollingCounter(10),
	}
}

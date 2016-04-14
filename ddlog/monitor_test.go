package ddlog

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
)

type monitorTestCase struct {
	c      *Config
	mclock *clock.Mock

	g *Generator
	m *Monitor
}

func (tc *monitorTestCase) Tick(dt time.Duration) {
	tc.mclock.Add(dt)
}

func newMonitorTestCase() *monitorTestCase {
	config := NewConfig()
	mclock := clock.NewMock()
	mclock.Set(time.Now())
	config.Mock(mclock)

	return &monitorTestCase{
		c:      config,
		mclock: mclock,
		g:      config.NewGenerator(),
		m:      config.NewMonitor(),
	}
}

func (tc *monitorTestCase) sendNMsg(msgChan chan *Message, count int) {
	for i := 0; i < count; i++ {
		msgChan <- tc.g.RandMsg()
	}
	time.Sleep(time.Millisecond) // allow proccessing
}

func TestMonitorAlert(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	tc := newMonitorTestCase()

	msgChan := make(chan *Message)
	ctx, cancel := context.WithCancel(context.Background())
	go tc.m.Start(ctx, msgChan)
	defer cancel()

	// Send exactly alert threshold
	tc.sendNMsg(msgChan, tc.c.AlertThreshold)
	a.Equal(tc.c.AlertThreshold, tc.m.WindowCount())
	a.Empty(tc.m.Alerts(), "Alert should happen at threshold+1 yet")

	// Send 1 more - should cause an alert.
	tc.sendNMsg(msgChan, 1)
	a.Len(tc.m.Alerts(), 1)
	alert := tc.m.Alerts()[0]
	a.False(alert.IsDone())

	// Send more events - should be only 1 event.
	tc.sendNMsg(msgChan, 1)
	a.Len(tc.m.Alerts(), 1)

	// Jump a window away - alert should stop
	tc.mclock.Add(tc.c.WindowSize * 2)
	a.Len(tc.m.Alerts(), 2)
	a.True(tc.m.Alerts()[1].IsDone())

	// create another alert.

	tc.sendNMsg(msgChan, tc.c.AlertThreshold+1)
	a.Len(tc.m.Alerts(), 3)
	a.False(tc.m.Alerts()[2].IsDone())
}

func TestMonitorStart(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	tc := newMonitorTestCase()

	m, g := tc.m, tc.g
	msgChan := make(chan *Message)
	ctx, cancel := context.WithCancel(context.Background())
	go m.Start(ctx, msgChan)
	defer cancel()

	// Visit a bunch of pages.
	visits := map[string]int{
		"/a": 1,
		"/b": 3,
		"/d": 9,
		"/x": 2,
	}

	totalVisitCount := 0
	for page, vc := range visits {
		totalVisitCount += vc
		for i := 0; i < vc; i++ {
			msgChan <- g.MsgWithPage(page)
		}
	}
	time.Sleep(time.Millisecond) // allow proccessing
	// check current state.
	top2 := m.TopK(2)
	a.Len(top2, 2)
	a.Equal("/d", top2[0].Name)
	a.Equal("/b", top2[1].Name)
	a.Equal(totalVisitCount, m.WindowCount())

	// Advance just before window expires
	tc.Tick(tc.c.WindowSize - 1*time.Second)
	a.Equal(totalVisitCount, m.WindowCount())

	tc.Tick(time.Second)
	a.Equal(0, m.WindowCount())
}

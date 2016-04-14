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
	config.Mock(mclock)

	return &monitorTestCase{
		c:      config,
		mclock: mclock,
		g:      config.NewGenerator(),
		m:      config.NewMonitor(),
	}
}

func (tc *monitorTestCase) sendAndProcess(msgChan chan *Message, count int) {
	ctx, cancel := context.WithCancel(context.Background())
	// create count messages
	go func() {
		for i := 0; i < count; i++ {
			msgChan <- tc.g.RandMsg()
		}
		cancel()
	}()

	tc.m.Start(ctx, msgChan)
}

func TestMonitorAlert(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	tc := newMonitorTestCase()

	msgChan := make(chan *Message)

	// Send exactly alert threshold
	tc.sendAndProcess(msgChan, tc.c.AlertThreshold)
	a.Equal(tc.c.AlertThreshold, tc.m.WindowCount())
	a.Empty(tc.m.Alerts(), "Alert should happen at threshold+1 yet")

	// Send 1 more - should cause an alert.
	tc.sendAndProcess(msgChan, 1)
	a.Len(tc.m.Alerts(), 1)
	alert := tc.m.Alerts()[0]
	a.False(alert.IsDone())

	// Send more events - should be only 1 event.
	tc.sendAndProcess(msgChan, 1)
	a.Len(tc.m.Alerts(), 1)

	// Jump a window away - alert should stop
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		tc.mclock.Add(tc.c.WindowSize * 2)
		cancel()
	}()
	tc.m.Start(ctx, msgChan)
	a.Len(tc.m.Alerts(), 2)
	a.True(tc.m.Alerts()[1].IsDone())

	// create another alert.

	tc.sendAndProcess(msgChan, tc.c.AlertThreshold+1)
	a.Len(tc.m.Alerts(), 3)
	a.False(tc.m.Alerts()[2].IsDone())
}

func TestMonitorStart(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	tc := newMonitorTestCase()

	m, g := tc.m, tc.g
	msgChan := make(chan *Message)

	// Visit a bunch of pages.
	visits := map[string]int{
		"/a": 1,
		"/b": 3,
		"/d": 9,
		"/x": 2,
	}
	visitCount := 0
	for _, vc := range visits {
		visitCount += vc
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for page, visitCount := range visits {
			for i := 0; i < visitCount; i++ {
				msgChan <- g.MsgWithPage(page)
			}
		}
		cancel()
	}()

	// Blocks until stopped
	m.Start(ctx, msgChan)

	// check current state.
	top2 := m.TopK(2)
	a.Len(top2, 2)
	a.Equal("/d", top2[0].Name)
	a.Equal("/b", top2[1].Name)
	a.Equal(visitCount, m.WindowCount())

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		// Advance 1m59s by seconds
		for i := 1; i < 120; i++ {
			tc.Tick(1 * time.Second)
		}
		cancel()
	}()

	m.Start(ctx, msgChan)
	a.Equal(visitCount, m.WindowCount())

	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		// Advance 1s more
		tc.Tick(1 * time.Second)
		cancel()
	}()
	m.Start(ctx, msgChan)
	// Old value's no longer exist because it's been 2 min.
	a.Equal(0, m.WindowCount())
}

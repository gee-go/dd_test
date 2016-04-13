package ddlog

import (
	"testing"
	"time"

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

func TestMonitorTick(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	tc := newMonitorTestCase()

	m, g := tc.m, tc.g
	msgChan := make(chan *Message)

	// Create 2 messages per second for 2 minutes
	total := 240
	tc.mclock.Set(time.Now().Round(time.Second))
	go func() {
		for i := 0; i < 120; i++ {
			tc.Tick(time.Second)
			msgChan <- g.RandMsg()
			msgChan <- g.RandMsg()
		}
		tc.Tick(time.Second)
		m.Stop()
	}()

	m.Start(msgChan)

	a.Equal(total, m.WindowCount())

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

	go func() {
		for page, visitCount := range visits {
			for i := 0; i < visitCount; i++ {
				msgChan <- g.MsgWithPage(page)
			}
		}
		m.Stop()
	}()

	// Blocks until stopped
	m.Start(msgChan)

	// check current state.
	top2 := m.TopK(2)
	a.Len(top2, 2)
	a.Equal("/d", top2[0].Name)
	a.Equal("/b", top2[1].Name)
	a.Equal(visitCount, m.WindowCount())

	go func() {
		// Advance 1m59s
		for i := 1; i < 120; i++ {
			tc.Tick(1 * time.Second)
		}
		m.Stop()
	}()
	m.Start(msgChan)
	a.Equal(visitCount, m.WindowCount())

	go func() {
		// Advance 1s more
		tc.Tick(1 * time.Second)
		m.Stop()
	}()
	m.Start(msgChan)
	// Old value's no longer exist because it's been 2 min.
	a.Equal(0, m.WindowCount())
}

package ddlog

import (
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gee-go/dd_test/ddlog/util"
)

type Monitor struct {
	rollingCount *util.CountRing
	clock        clock.Clock
}

func NewMonitor() *Monitor {
	return &Monitor{
		rollingCount: util.NewCountRing(time.Second, 240),
		clock:        clock.New(),
	}
}

func (m *Monitor) Start(msgChan <-chan *Message) {
	intervalTick := m.clock.Tick(time.Second)

	for {
		select {
		case msg := <-msgChan:
			// TODO - check for dropped messages
			m.rollingCount.Inc(msg.Time, 1)
		case <-intervalTick:
			m.rollingCount.Tick()
		}
	}
}

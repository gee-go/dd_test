package ddlog

import (
	"sync"

	"github.com/gee-go/ddlog/ddlog/util"
)

// Monitor monitors a message channel and aggregates stats.
type Monitor struct {
	// mutable structures
	mu           sync.RWMutex
	rollingCount *util.CountRing
	pageCount    *MetricCounter

	*Config
}

func NewMonitor(c *Config) *Monitor {
	return &Monitor{
		rollingCount: util.NewCountRing(c.AggInterval, c.numWindowsKept()),
		pageCount:    NewMetricCounter(),
		Config:       c,
	}
}

// TopK returns a sorted list of the top k visited pages since the program started.
func (m *Monitor) TopK(k int) []*PageCount {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.pageCount.TopK(k)
}

// Start collecting metrics on the messages
func (m *Monitor) Start(msgChan <-chan *Message) {
	intervalTick := m.clock.Tick(m.Config.AggInterval)

	for {
		select {
		case msg := <-msgChan:
			m.mu.Lock()
			m.rollingCount.Inc(msg.Time, 1)
			m.mu.Unlock()
		case <-intervalTick:
			// updates the rolling count when no messages have arrived in the last interval.
			m.mu.Lock()
			m.rollingCount.Tick()
			m.mu.Unlock()
		}
	}
}

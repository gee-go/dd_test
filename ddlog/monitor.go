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

	c    *Config
	stop chan bool
}

func NewMonitor(c *Config) *Monitor {
	return &Monitor{
		rollingCount: util.NewCountRing(c.AggInterval, c.numWindowsKept(), c.clock),
		pageCount:    NewMetricCounter(),
		c:            c,
		stop:         make(chan bool),
	}
}

func (m *Monitor) Stop() {
	m.stop <- true
}

func (m *Monitor) WindowCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rollingCount.Sum()
}

// TopK returns a sorted list of the top k visited pages since the program started.
func (m *Monitor) TopK(k int) []*PageCount {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.pageCount.TopK(k)
}

// Start collecting metrics on the messages
func (m *Monitor) Start(msgChan <-chan *Message) {
	intervalTicker := m.c.clock.Ticker(m.c.AggInterval)

	defer intervalTicker.Stop()
	for {
		select {
		case msg := <-msgChan:
			m.mu.Lock()
			m.rollingCount.Inc(msg.Time, 1)
			m.pageCount.Inc(msg)
			m.mu.Unlock()
		case <-intervalTicker.C:
			// updates the rolling count when no messages have arrived in the last interval.
			m.mu.Lock()
			m.rollingCount.Tick()
			m.mu.Unlock()
		case <-m.stop:
			return
		}

	}
}

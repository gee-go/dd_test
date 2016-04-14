package ddlog

import (
	"sync"

	"golang.org/x/net/context"

	"github.com/gee-go/ddlog/ddlog/util"
)

// Monitor monitors a message channel and aggregates stats.
type Monitor struct {
	// mutable structures
	mu           sync.RWMutex
	rollingCount *util.CountRing
	pageCount    *MetricCounter
	alert        *Alert

	alerts []*Alert // TODO - this grows forever.

	c *Config

	ctx    context.Context
	cancel context.CancelFunc
}

func NewMonitor(c *Config) *Monitor {
	return &Monitor{
		rollingCount: util.NewCountRing(c.AggInterval, c.numWindowsKept(), c.clock),
		pageCount:    NewMetricCounter(),
		c:            c,
	}
}

func (m *Monitor) Config() *Config {
	return m.c
}

func (m *Monitor) Spark() []float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rollingCount.Spark()
}

func (m *Monitor) Stop() {
	m.cancel()
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

func (m *Monitor) Alerts() []*Alert {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.alerts
}

func (m *Monitor) checkAlert() {
	count := m.rollingCount.Sum()
	isAlertMode := count > m.c.AlertThreshold

	// already have an alert.
	if isAlertMode && m.alert == nil {
		// New alert
		m.alert = &Alert{Start: m.c.clock.Now(), Count: count}
		m.alerts = append(m.alerts, m.alert.Copy())
	}

	if m.alert != nil && !isAlertMode {
		// stop alert.
		m.alert.Complete(m.c.clock.Now())
		m.alerts = append(m.alerts, m.alert.Copy())
		m.alert = nil
	}
}

// Start collecting metrics on the messages
func (m *Monitor) Start(msgChan <-chan *Message) {
	m.ctx, m.cancel = context.WithCancel(context.Background())

	go func() {
		intervalTicker := m.c.clock.Ticker(m.c.AggInterval)
		defer intervalTicker.Stop()
		for {
			select {
			case <-intervalTicker.C:
				// updates the rolling count when no messages have arrived in the last interval.
				m.mu.Lock()
				m.rollingCount.Tick()
				m.checkAlert()
				m.mu.Unlock()
			case <-m.ctx.Done():
				return
			}
		}
	}()

	for {
		select {
		case msg := <-msgChan:
			m.mu.Lock()
			m.rollingCount.Inc(msg.Time, 1)
			m.pageCount.Inc(msg)
			m.checkAlert()
			m.mu.Unlock()
		case <-m.ctx.Done():
			return
		}

	}
}

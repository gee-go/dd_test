package ddlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/eapache/channels"
	"github.com/gee-go/ddlog/ddlog/util"
)

type Alert struct {
	Start time.Time
	End   time.Time
	Count int
}

func (a *Alert) Send(ch chan *Alert) {
	select {
	case ch <- a.Copy():
	// no block
	default:
		return
	}
}
func (a *Alert) String() string {
	if a.IsDone() {
		return fmt.Sprintf("[Alert Done] at %s duration=%s\n", a.End, a.End.Sub(a.Start))
	}
	return fmt.Sprintf("High traffic generated an alert - hits = %v, triggered at %s", a.Count, a.Start)
}

func (a *Alert) Copy() *Alert {
	return &Alert{Start: a.Start, End: a.End, Count: a.Count}
}

func (a *Alert) Complete(at time.Time) {
	a.End = at
}

func (a *Alert) IsDone() bool {
	return !a.End.IsZero()
}

// Monitor monitors a message channel and aggregates stats.
type Monitor struct {
	// mutable structures
	mu           sync.RWMutex
	rollingCount *util.CountRing
	pageCount    *MetricCounter
	alert        *Alert

	alertChan *channels.InfiniteChannel

	c    *Config
	stop chan bool
}

func NewMonitor(c *Config) *Monitor {
	return &Monitor{
		rollingCount: util.NewCountRing(c.AggInterval, c.numWindowsKept(), c.clock),
		pageCount:    NewMetricCounter(),
		c:            c,
		stop:         make(chan bool),
		alertChan:    channels.NewInfiniteChannel(),
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

func (m *Monitor) AlertChan() <-chan interface{} {
	return m.alertChan.Out()
}

func (m *Monitor) checkAlert() {
	count := m.rollingCount.Sum()
	isAlertMode := count > m.c.AlertThreshold

	// already have an alert.
	if isAlertMode && m.alert == nil {
		// New alert
		m.alert = &Alert{Start: m.c.clock.Now(), Count: count}
		m.alertChan.In() <- m.alert.Copy()
	}

	if m.alert != nil && !isAlertMode {
		// stop alert.
		m.alert.Complete(m.c.clock.Now())
		m.alertChan.In() <- m.alert.Copy()
		m.alert = nil
	}
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
			m.checkAlert()
			m.mu.Unlock()
		case <-intervalTicker.C:
			// updates the rolling count when no messages have arrived in the last interval.
			m.mu.Lock()
			m.rollingCount.Tick()
			m.checkAlert()
			m.mu.Unlock()
		case <-m.stop:
			return
		}

	}
}

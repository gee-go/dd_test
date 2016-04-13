package ddlog

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

type Alert struct {
	Start time.Time
	End   time.Time
	Count int
	Done  bool
}

func (a *Alert) Copy() *Alert {
	if a == nil {
		return nil
	}
	return &Alert{
		Start: a.Start,
		End:   a.End,
		Done:  a.Done,
	}
}

type MetricEvent struct {
	Alert    *Alert
	TopPages []*PageCount
}

type MetricEventHandler func(e *MetricEvent)

type MetricStore struct {
	mu sync.RWMutex

	config        *Config
	clock         clock.Clock
	allTimeBucket *MetricBucket // metrics from the beginning of this process

	fastBucketRing *metricBucketRing // a ring of buckets represending intervals of 10s for the past 2 mins.
	msgChan        <-chan *Message
	fastTicker     *clock.Ticker
	metricEventFn  MetricEventHandler

	alert *Alert
}

func NewMetricStore(config *Config) *MetricStore {
	return &MetricStore{
		config:         config,
		clock:          clock.New(),
		allTimeBucket:  NewMetricBucket(),
		fastBucketRing: newMetricBucketRing(config.FastTickDuration, int(config.AlertDuration/config.FastTickDuration)),
	}
}

func (ms *MetricStore) step() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.fastBucketRing.Step()

	// alerts
	alertBucket := ms.fastBucketRing.Merged()
	var a *Alert

	// new alert
	if alertBucket.Count >= ms.config.AlertThreshold && ms.alert == nil {
		ms.alert = &Alert{
			Start: time.Now(),
			Count: alertBucket.Count,
		}
		a = ms.alert.Copy()
	}

	// alert done
	if alertBucket.Count < ms.config.AlertThreshold && ms.alert != nil {
		a = ms.alert.Copy()
		a.Done = true
		a.End = time.Now()
		ms.alert = nil
	}

	evt := &MetricEvent{
		Alert:    a,
		TopPages: ms.allTimeBucket.TopK(5),
	}

	go ms.metricEventFn(evt)

}

func (ms *MetricStore) Start(msgChan <-chan *Message, handler MetricEventHandler) {
	ms.msgChan = msgChan
	ms.fastTicker = ms.clock.Ticker(ms.config.FastTickDuration)
	ms.metricEventFn = handler

	go func() {
		for range ms.fastTicker.C {
			ms.step()
		}
	}()

	go func() {
		for m := range msgChan {
			ms.Add(m)
		}
		ms.fastTicker.Stop()
	}()

}

func (ms *MetricStore) Add(m *Message) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.allTimeBucket.Add(m)
	ms.fastBucketRing.Add(m)
}

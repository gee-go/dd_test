package ddlog

import (
	"sync"

	"github.com/benbjohnson/clock"
)

type MetricEvent struct {
	AlertBucket *MetricBucket
	TopPages    []*PageCount
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

	go ms.metricEventFn(&MetricEvent{
		AlertBucket: ms.fastBucketRing.Merged(),
		TopPages:    ms.allTimeBucket.TopK(5),
	})
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

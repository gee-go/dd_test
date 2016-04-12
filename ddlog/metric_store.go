package ddlog

import (
	"fmt"
	"sync"
	"time"
)

type MetricStore struct {
	mu sync.RWMutex

	*MetricBucket
}

func NewMetricStore() *MetricStore {
	return &MetricStore{}
}

func (s *MetricStore) Print() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("Metric Print:")
	for _, p := range s.TopK(5) {
		fmt.Println(p, s.CountByPage[p])
	}

}

// Start calculating aggregate stats from the msgChan, messages are rolled up
// in aggInterval buckets (default 1s).
func (s *MetricStore) Start(msgChan <-chan *Message, aggInterval ...time.Duration) {
	aggI := time.Second * 1
	if len(aggInterval) == 1 {
		aggI = aggInterval[0]
	}

	bChan := MsgChanToBucket(msgChan, aggI)

	for b := range bChan {
		s.merge(b)
	}
}

func (s *MetricStore) merge(b *MetricBucket) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.MetricBucket == nil {
		s.MetricBucket = b
		return
	}

	s.Count += b.Count

	for page, count := range b.CountByPage {
		s.CountByPage[page] += count
	}
}

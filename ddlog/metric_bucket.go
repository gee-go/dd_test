package ddlog

import (
	"fmt"
	"time"

	"github.com/gee-go/dd_test/ddlog/util"
)

// MetricBucket stores aggregated events over a period of time.
type MetricBucket struct {
	StartTime time.Time

	Count       int
	CountByPage map[string]int
}

// Add the message to the aggregate. Not thread safe!
func (b *MetricBucket) Add(m *Message) {
	b.Count++
	name := m.EventName()
	if _, found := b.CountByPage[name]; !found {
		b.CountByPage[name] = 1
	} else {
		b.CountByPage[name]++
	}
}

func (b *MetricBucket) TopK(k int) []string {
	return util.TopK(b.CountByPage, k)
}

func NewMetricBucket(mt time.Time) *MetricBucket {
	return &MetricBucket{
		CountByPage: make(map[string]int),
		StartTime:   mt,
	}
}

func MsgChanToBucket(msgChan <-chan *Message, aggInterval time.Duration) chan *MetricBucket {
	out := make(chan *MetricBucket)

	go func() {
		var b *MetricBucket

		for m := range msgChan {
			mt := m.Time.Truncate(aggInterval)

			if b == nil {
				// first bucket
				b = NewMetricBucket(mt)
			} else if mt.After(b.StartTime) {
				// new bucket
				out <- b
				b = NewMetricBucket(mt)
			}

			if mt.Equal(b.StartTime) {
				// same bucket
				b.Add(m)
			} else {
				// TODO - Handle case where log timestamps are not constantly increasing.
				fmt.Println("Old line")
			}

		}
		close(out)
	}()

	return out
}

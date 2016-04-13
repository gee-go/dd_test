package ddlog

import (
	"time"

	"github.com/gee-go/ddlog/ddlog/util"
)

type PageCount struct {
	Name  string
	Count int
}

// MetricBucket stores aggregated events over a period of time.
type MetricBucket struct {
	Duration time.Duration

	Count       int
	CountByPage map[string]int
}

// Copy creates a copy of this bucket
func (b *MetricBucket) Copy() *MetricBucket {
	out := NewMetricBucket()

	out.Merge(b)

	return out
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

// Merge modifies this bucket
func (b *MetricBucket) Merge(in *MetricBucket) {
	// TODO StartTime
	b.Count += in.Count
	b.Duration += in.Duration

	for page, count := range in.CountByPage {
		b.CountByPage[page] += count
	}
}

func (b *MetricBucket) TopK(k int) []*PageCount {
	keys := util.TopK(b.CountByPage, k)

	var out []*PageCount
	for _, page := range keys {
		out = append(out, &PageCount{
			Name:  page,
			Count: b.CountByPage[page],
		})
	}

	return out
}

func NewMetricBucket() *MetricBucket {
	return &MetricBucket{
		CountByPage: make(map[string]int),
	}
}

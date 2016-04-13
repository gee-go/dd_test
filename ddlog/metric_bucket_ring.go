package ddlog

import (
	"fmt"
	"time"
)

// a ring of recent metric buckets. Not thread safe!
type metricBucketRing struct {
	data          []*MetricBucket
	current       int
	batchInterval time.Duration
	retainCount   int
}

func newMetricBucketRing(batchInterval time.Duration, retainCount int) *metricBucketRing {
	fmt.Println(retainCount)
	rc := &metricBucketRing{
		data:          make([]*MetricBucket, retainCount),
		batchInterval: batchInterval,
		retainCount:   retainCount,
	}
	rc.resetBucket(0)
	return rc
}

func (w *metricBucketRing) resetBucket(i int) {
	b := NewMetricBucket()
	b.Duration = w.batchInterval
	w.data[i] = b
}

func (w *metricBucketRing) cur() *MetricBucket {
	return w.data[w.current]
}

// Merged returns a new bucket that is the sum of every bucket in the ring.
func (w *metricBucketRing) Merged() *MetricBucket {
	out := NewMetricBucket()
	for _, b := range w.data {
		if b == nil {
			continue
		}

		out.Merge(b)
	}

	return out
}

// Step moves on to the next bucket, and clears/returns the oldest bucket.
func (w *metricBucketRing) Step() *MetricBucket {
	prev := w.cur()

	w.current = (w.current + 1) % len(w.data)
	w.resetBucket(w.current)
	return prev.Copy()
}

// Add a message to the current bucket
func (w *metricBucketRing) Add(m *Message) {
	w.data[w.current].Add(m)
}

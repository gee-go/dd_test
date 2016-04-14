package ddlog

import "github.com/gee-go/ddlog/ddlog/util"

type PageCount struct {
	Name         string
	Count        int
	CountPercent float64
	Bytes        uint64
}

// MetricCounter tracks a total count and per page count.
type MetricCounter struct {
	totalCount int

	countByPage map[string]int
	bytesByPage map[string]uint64
}

func NewMetricCounter() *MetricCounter {
	return &MetricCounter{
		countByPage: make(map[string]int),
		bytesByPage: make(map[string]uint64),
	}
}

func (b *MetricCounter) Inc(msg *Message) {
	page := msg.EventName()

	// Count
	b.totalCount += 1
	if _, found := b.countByPage[page]; !found {
		b.countByPage[page] = 1
	} else {
		b.countByPage[page] += 1
	}

	// Bytes
	if _, found := b.bytesByPage[page]; !found {
		b.bytesByPage[page] = uint64(msg.Size)
	} else {
		b.bytesByPage[page] += uint64(msg.Size)
	}
}

func (b *MetricCounter) TopK(k int) []*PageCount {
	keys := util.TopK(b.countByPage, k)

	var out []*PageCount

	for _, page := range keys {
		count := b.countByPage[page]
		out = append(out, &PageCount{
			Name:         page,
			Count:        count,
			CountPercent: float64(count) / float64(b.totalCount),
			Bytes:        b.bytesByPage[page],
		})
	}

	return out
}

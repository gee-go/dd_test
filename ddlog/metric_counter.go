package ddlog

import "github.com/gee-go/ddlog/ddlog/util"

type PageCount struct {
	Name  string
	Count int
}

// MetricCounter tracks a total count and per page count.
type MetricCounter struct {
	totalCount  int
	countByPage map[string]int
}

func NewMetricCounter() *MetricCounter {
	return &MetricCounter{
		countByPage: make(map[string]int),
	}
}

func (b *MetricCounter) Inc(page string, by int) {
	b.totalCount += by
	if _, found := b.countByPage[page]; !found {
		b.countByPage[page] = by
	} else {
		b.countByPage[page] += by
	}
}

func (b *MetricCounter) TopK(k int) []*PageCount {
	keys := util.TopK(b.countByPage, k)

	var out []*PageCount
	for _, page := range keys {
		out = append(out, &PageCount{
			Name:  page,
			Count: b.countByPage[page],
		})
	}

	return out
}

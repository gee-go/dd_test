package lscan

import (
	"log"
	"os"
	"time"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/rcrowley/go-metrics"
)

type MetricStore struct {
	count metrics.Counter
}

func NewMetric() *MetricStore {
	ms := &MetricStore{}
	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	return ms
}

func (ms *MetricStore) HandleMsg(m *lparse.Message) {
	metrics.GetOrRegisterCounter("cc", metrics.DefaultRegistry).Inc(1)
	// ms.sink.IncrCounter([]string{"count"}, 1)
}

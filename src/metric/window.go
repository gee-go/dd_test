package metric

import (
	"sync"
	"time"

	"github.com/gee-go/dd_test/src/lparse"
)

type Window struct {
	data [][]*lparse.Message

	mu      sync.Mutex
	current int
}

func NewWindow(size int) *Window {
	return &Window{
		data: make([][]*lparse.Message, size),
	}
}

func (w *Window) Step() {
	// move current, and clear old
	w.current = (w.current + 1) % len(w.data)
	w.data[w.current] = w.data[w.current][:0]
}

func (w *Window) Insert(m *lparse.Message) {
	w.data[w.current] = append(w.data[w.current], m)
}

type MetricStore struct {
	windowDuration time.Duration
	window         *Window
}

func NewMetric() *MetricStore {
	return &MetricStore{
		windowDuration: time.Second * 2,
		window:         NewWindow(3),
	}
}

// func (ms *MetricStore) Start() {
// 	for range time.Tick(ms.windowDuration) {
// 		fmt.Println("tick", len(ms.window.Current()))

// 		ms.window.Step()
// 	}
// }

func (ms *MetricStore) HandleMsg(m *lparse.Message) {
	ms.window.Insert(m)
}

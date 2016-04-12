package metric

import (
	"container/heap"
	"sync"

	"github.com/gee-go/dd_test/ddlog"
)

type Store struct {
	idToPage   map[string]*Page
	totalCount int
	mu         sync.RWMutex
}

func (s *Store) HandleMsg(m *ddlog.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := m.EventName()
	s.totalCount++
	s.incPageStatus(name, m.Status)
}

func (s *Store) incPageStatus(id string, status int) {
	p, found := s.idToPage[id]
	if !found {
		p = NewPage(id)
		s.idToPage[id] = p
	}
	p.IncStatus(status)
}

func (s *Store) TopK(k int) []Page {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h := &Pages{}

	i := 0
	for _, p := range s.idToPage {
		// put first k in heap
		if i < k {
			heap.Push(h, p)
			i++
		} else if p.Total > (*h)[0].Total {
			(*h)[0] = *p
			heap.Fix(h, 0)
		}
	}

	out := make([]Page, k)
	for i = 0; i < k; i++ {
		out[i] = heap.Pop(h).(Page)
	}

	return out
}

func New() *Store {
	return &Store{
		idToPage: make(map[string]*Page),
	}
}

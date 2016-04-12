package metric

import (
	"container/heap"

	"github.com/gee-go/dd_test/src/lparse"
)

type CountByStatus struct {
	c1xx, c2xx, c3xx, c4xx, c5xx int

	total int
}

type Store struct {
	idToPage   map[string]*Page
	totalCount int
}

func (s *Store) HandleMsg(m *lparse.Message) {
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
	p.Count++
	// p.Count.IncStatus(status)
}

func (s *Store) TopK(k int) []*Page {
	h := &Pages{}

	i := 0
	for _, p := range s.idToPage {
		// put first k in heap
		if i < k {
			heap.Push(h, p)
			i++
		} else if p.Count > (*h)[0].Count {
			(*h)[0] = p
			heap.Fix(h, 0)
		}
	}

	out := make([]*Page, k)
	for i = 0; i < k; i++ {
		out[i] = heap.Pop(h).(*Page)
	}

	return out
}

func New() *Store {
	return &Store{
		idToPage: make(map[string]*Page),
	}
}

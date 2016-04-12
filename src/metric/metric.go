package metric

import "container/heap"

type CountByStatus struct {
	c1xx, c2xx, c3xx, c4xx, c5xx int

	total int
}

func (cs *CountByStatus) IncStatus(status int) {
	cs.total++

	switch {
	case status < 200:
		cs.c1xx++
	case status < 300:
		cs.c2xx++
	case status < 400:
		cs.c3xx++
	case status < 500:
		cs.c4xx++
	default:
		cs.c5xx++
	}

}

type Page struct {
	ID string // e.g. /store

	Count int
}

func NewPage(id string) *Page {
	return &Page{
		ID: id,
	}
}

type Store struct {
	idToPage map[string]*Page
}

func (s *Store) IncPageStatus(id string, status int) {
	p, found := s.idToPage[id]
	if !found {
		p = NewPage(id)
		s.idToPage[id] = p
	}

	p.Count.IncStatus(status)
}

type IntHeap []*Page

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Page))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (s *Store) TopK(k int) []*Page {
	h := &IntHeap{}

	i := 0
	for _, p := range s.idToPage {
		if i < k {
			heap.Push(h, p)
		} else if p.Count > h[0].Count {
			h[0] = p
			heap.Fix(h, 0)
		}

		i++
	}
}

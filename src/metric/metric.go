package metric

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

	Count *CountByStatus
}

func NewPage(id string) *Page {
	return &Page{
		ID:    id,
		Count: &CountByStatus{},
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

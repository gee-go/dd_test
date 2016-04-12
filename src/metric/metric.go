package metric

type CountByStatus struct {
	Count1xx, Count2xx, Count3xx, Count4xx, Count5xx int
	Total                                            int
}

func (cs *CountByStatus) IncStatus(status int) {
	cs.Total++
	switch {
	case status < 200:
		cs.Count1xx++
	case status < 300:
		cs.Count2xx++
	case status < 400:
		cs.Count3xx++
	case status < 500:
		cs.Count4xx++
	default:
		cs.Count5xx++
	}

}

type Page struct {
	ID string // e.g. /store

	CountByStatus
}

func NewPage(id string) *Page {
	return &Page{
		ID: id,
	}
}

type Pages []Page

// Len - sort.Interface
func (h Pages) Len() int { return len(h) }

// Swap - sort.Interface
func (h Pages) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h Pages) Less(i, j int) bool { return h[i].Total < h[j].Total }

func (h *Pages) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Page))
}

func (h *Pages) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

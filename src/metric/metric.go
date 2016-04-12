package metric

type Page struct {
	ID string // e.g. /store

	Count int // Event Count.
}

func NewPage(id string) *Page {
	return &Page{
		ID: id,
	}
}

type Pages []*Page

// Len - sort.Interface
func (h Pages) Len() int { return len(h) }

// Swap - sort.Interface
func (h Pages) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h Pages) Less(i, j int) bool { return h[i].Count < h[j].Count }

func (h *Pages) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Page))
}

func (h *Pages) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

package util

import (
	"container/heap"
	"sort"
)

type countGroupHeap struct {
	data map[string]int
	keys []string
}

func (h *countGroupHeap) Len() int {
	// using keys is important, only want k items in heap
	return len(h.keys)
}

func (h *countGroupHeap) get(i int) int {
	return h.data[h.keys[i]]
}

func (h *countGroupHeap) Less(i, j int) bool {
	// sort by count then key
	iv, jv := h.get(i), h.get(j)

	if iv == jv {
		return h.keys[i] < h.keys[j]
	}

	return iv < jv
}
func (h *countGroupHeap) Swap(i, j int) { h.keys[i], h.keys[j] = h.keys[j], h.keys[i] }

func (h *countGroupHeap) Push(x interface{}) {
	h.keys = append(h.keys, x.(string))
}

func (h *countGroupHeap) Pop() interface{} {
	var x string
	x, h.keys = h.keys[len(h.keys)-1], h.keys[:len(h.keys)-1]
	return x
}

// TopK keys by their value from data. O(n log k)
func TopK(data map[string]int, k int) []string {
	if k == 0 {
		var o []string
		return o
	}

	i := 0
	h := &countGroupHeap{
		data: data,
	}

	for page, count := range data {
		if k > i {
			heap.Push(h, page)
			i++
		} else if count > h.get(0) {
			h.keys[0] = page
			heap.Fix(h, 0)
		}
	}

	sort.Sort(sort.Reverse(h))

	return h.keys
}

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTopK(t *testing.T) {
	a := require.New(t)
	data := map[string]int{
		"a": 3,
		"b": 10,
		"c": 1,
		"d": 2,
		"e": 5,
	}
	order := []string{"b", "e", "a", "d", "c"}

	a.Equal(order, TopK(data, len(data)))
	a.Equal(order, TopK(data, len(data)+10))
	a.Equal([]string{"b", "e", "a"}, TopK(data, 3))
	delete(data, "a")
	a.Equal([]string{"b", "e", "d"}, TopK(data, 3))
}

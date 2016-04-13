package ddlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMetricBucketRing(t *testing.T) {
	a := require.New(t)
	g := NewGenerator(nil)
	batchInterval := 5 * time.Second
	retainCount := 10

	ring := newMetricBucketRing(batchInterval, retainCount)

	// test empty
	{
		b := ring.Merged()
		a.Equal(b.Count, 0)
		a.Len(b.TopK(10), 0)
		a.Equal(b.Duration, batchInterval)
	}

	{
		ring.Add(g.MsgWithPage("a"))
		b := ring.Merged()
		a.Equal(b.Count, 1)
		a.Len(b.TopK(10), 1)
		a.Equal(b.Duration, batchInterval)
		a.Equal(b.CountByPage["/a"], 1)
	}

	{
		ring.Add(g.MsgWithPage("a"))
		b := ring.Merged()
		a.Equal(ring.current, 0)
		a.Equal(b.Count, 2)
		a.Len(b.TopK(10), 1)
		a.Equal(b.Duration, batchInterval)
		a.Equal(b.CountByPage["/a"], 2)
	}

	for i := 1; i < retainCount; i++ {
		ring.Step()
		a.Equal(ring.current, i)

		b := ring.Merged()
		a.Equal(b.Count, 2)

		a.Equal(b.Duration, time.Duration(i+1)*batchInterval)
		a.Equal(b.CountByPage["/a"], 2)
	}

	// old data deleted
	ring.Step()
	b := ring.Merged()
	a.Equal(b.Count, 0)

	// a.Equal(ring.Merged().CountByPage["/a"], 2)
}

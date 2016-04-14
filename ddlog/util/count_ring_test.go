package util

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/require"
)

func TestCountRing(t *testing.T) {
	a := require.New(t)
	ring := NewCountRing(time.Second, 3)
	start := time.Now().Round(ring.dtInterval)

	ring.Inc(start, 1)
	a.Equal([]int{1, 0, 0}, ring.ring)
	a.Equal(0, ring.i)

	ring.Inc(start.Add(time.Millisecond*500), 1)
	a.Equal([]int{2, 0, 0}, ring.ring)
	a.Equal(0, ring.i)

	ring.Inc(start.Add(time.Second), 1)
	a.Equal([]int{2, 1, 0}, ring.ring)
	a.Equal(1, ring.i)
	a.Equal(3, ring.Sum())
}

func TestCountRingSpark(t *testing.T) {
	a := require.New(t)
	ring := NewCountRing(time.Second, 3)
	start := time.Now().Round(ring.dtInterval)

	incAt := func(step int, by int) bool {
		t := start.Add(time.Second * time.Duration(step))
		return ring.Inc(t, by)
	}

	incAt(0, 1)
	incAt(1, 2)
	incAt(2, 3)

	a.Equal([]int{1, 2, 3}, ring.ring)
	a.Equal(2, ring.i)
	a.Equal([]float64{1, 2, 3}, ring.Spark())

	ring.i = 1
	a.Equal([]float64{3, 1, 2}, ring.Spark())

	ring.i = 0
	a.Equal([]float64{2, 3, 1}, ring.Spark())
}

func TestCountRingOldMessages(t *testing.T) {
	a := require.New(t)
	ring := NewCountRing(time.Second, 3)
	start := time.Now().Round(ring.dtInterval)

	incAt := func(step int, by int) bool {
		t := start.Add(time.Second * time.Duration(step))
		return ring.Inc(t, by)
	}

	incAt(0, 1)
	incAt(1, 2)
	incAt(2, 3)

	a.Equal([]int{1, 2, 3}, ring.ring)
	a.Equal(2, ring.i)

	// Ago
	a.Equal(3, ring.Ago(0))
	a.Equal(2, ring.Ago(1))
	a.Equal(1, ring.Ago(2))
	a.Equal(-1, ring.Ago(3))

	a.True(incAt(0, 1))
	a.True(incAt(1, 2))
	a.True(incAt(2, 3))
	a.Equal([]int{2, 4, 6}, ring.ring)
	a.Equal(2, ring.i)

	a.True(incAt(3, 3))
	a.Equal([]int{3, 4, 6}, ring.ring)
	a.Equal(0, ring.i)

	a.True(incAt(2, 3))
	a.Equal([]int{3, 4, 9}, ring.ring)
	a.Equal(0, ring.i)

	// too far ago
	a.False(incAt(0, 3))
	a.Equal([]int{3, 4, 9}, ring.ring)
	a.Equal(0, ring.i)
}

func TestCountRingTick(t *testing.T) {
	a := require.New(t)
	mclock := clock.NewMock()
	ring := NewCountRing(time.Second, 3, mclock)

	start := ring.clock.Now().Round(ring.dtInterval)
	ring.Inc(start, 1)
	a.Equal([]int{1, 0, 0}, ring.ring)
	a.Equal(0, ring.i)

	ring.Tick()
	a.Equal(0, ring.i)
	mclock.Add(time.Second)
	ring.Tick()
	a.Equal(1, ring.i)

}

func TestCountRingSkip(t *testing.T) {
	a := require.New(t)
	ring := NewCountRing(time.Second, 5)
	start := time.Now().Round(ring.dtInterval)

	ring.Inc(start, 1)
	a.Equal([]int{1, 0, 0, 0, 0}, ring.ring)
	a.Equal(0, ring.i)

	ring.Inc(start.Add(2*time.Second), 2)
	a.Equal([]int{1, 0, 2, 0, 0}, ring.ring)
	a.Equal(2, ring.i)

	ring.Inc(start.Add(5*time.Second), 3)
	a.Equal([]int{3, 0, 2, 0, 0}, ring.ring)
	a.Equal(0, ring.i)
}

func BenchmarkCountRing(b *testing.B) {
	ring := NewCountRing(time.Second, 60)
	start := time.Now().Round(ring.dtInterval)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ring.Inc(start.Add(time.Millisecond*time.Duration(i)), 1)
	}
}

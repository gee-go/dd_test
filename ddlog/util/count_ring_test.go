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

func TestCountRingTick(t *testing.T) {
	a := require.New(t)
	ring := NewCountRing(time.Second, 3)
	mclock := clock.NewMock()
	ring.clock = mclock
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

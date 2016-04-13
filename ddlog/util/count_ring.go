package util

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

type CountRing struct {
	mu         sync.RWMutex
	dtInterval time.Duration

	i    int
	ring []int

	current     time.Time // the truncated time of the current bucket.
	lastMsgTime time.Time // the most recent message time.
	clock       clock.Clock
}

func NewCountRing(dtInterval time.Duration, size int) *CountRing {
	return &CountRing{
		dtInterval: dtInterval,
		ring:       make([]int, size),
		clock:      clock.New(),
	}
}

func (r *CountRing) advance(by int) {
	r.i = (r.i + by) % len(r.ring)
	r.ring[r.i] = 0
}

// Tick advances current even in the absence of messages.
// Try to call once per dtInterval.
func (r *CountRing) Tick() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock.Now()
	dt := now.Sub(r.lastMsgTime)
	if dt >= r.dtInterval {
		r.lastMsgTime = now
		r.advance(int(dt / r.dtInterval))
	}
}

// Sum of the data.
func (r *CountRing) Sum() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s := 0
	for _, v := range r.ring {
		s += v
	}
	return s
}

// Inc the bucket corresponding to time `at`.
// Discard if the time is before the current time.
func (r *CountRing) Inc(at time.Time, by int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	t := at.Truncate(r.dtInterval)

	if r.current.IsZero() {
		r.current = t
	}

	r.lastMsgTime = r.clock.Now()
	if t.Equal(r.current) {
		r.ring[r.i] += by
	} else if t.After(r.current) {
		dt := int(t.Sub(r.current) / r.dtInterval)
		r.current = t
		r.advance(dt)
		r.ring[r.i] += by
	} else {
		return false
	}

	return true
}

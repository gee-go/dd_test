package randutil

import (
	"math/rand"
	"net"
	"time"

	"golang.org/x/exp/utf8string"
)

var (
	R      = New()
	alphaC = utf8string.NewString("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
)

type Rand struct {
	Rand *rand.Rand
}

func New() *Rand {

	return Quick(rand.New(&lockedSource{src: rand.NewSource(time.Now().UnixNano())}))
}

func Quick(r *rand.Rand) *Rand {
	return &Rand{r}
}

// Bytes generates n bytes
func (r *Rand) Bytes(n int) []byte {
	b := make([]byte, n)
	r.Rand.Read(b) // its ok to ignore err https://golang.org/pkg/math/rand/#Read
	return b
}

// IPv4 returns a random IP address
func (r *Rand) IPv4() net.IP {
	return net.IP(r.Bytes(4))
}

// IntRange returns a non-negative pseudo-random number in [min,max)
func (r *Rand) IntRange(min, max int) int {
	return min + r.Rand.Intn(max-min)
}

// Alpha returns a random string of len(n) with azAZ chars.
func (r *Rand) Alpha(n int) string {
	return r.string(n, alphaC)
}

func (r *Rand) string(n int, choices *utf8string.String) string {
	l := choices.RuneCount()

	out := make([]rune, n)
	for i := 0; i < n; i++ {
		out[i] = choices.At(r.Rand.Intn(l))
	}
	return string(out)
}

// SelectString returns a random string from the given choices.
// Returns "" if no choices are given
func (r *Rand) SelectString(choices ...string) string {
	if len(choices) == 0 {
		return ""
	}

	return choices[r.Rand.Intn(len(choices))]

}

func (r *Rand) SelectInt(choices ...int) int {
	if len(choices) == 0 {
		return 0
	}

	return choices[r.Rand.Intn(len(choices))]

}

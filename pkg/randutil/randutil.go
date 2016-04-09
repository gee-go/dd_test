package randutil

import (
	"bytes"
	"math/rand"
	"net"
	"time"
	"unicode/utf8"
)

const (
	alpha   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	letters = alpha + "日a本b語ç日ð本Ê語þ日¥本¼語i日©日a本b語ç日ð本Ê語þ日¥本¼語i日©日a本b語ç日ð本Ê語þ日¥本¼語i日©"
)

type Rand struct {
	rand *rand.Rand
}

func New() *Rand {
	return FromSeed(time.Now().UnixNano())
}

func FromSeed(seed int64) *Rand {
	return Quick(rand.New(rand.NewSource(seed)))
}

func Quick(r *rand.Rand) *Rand {
	return &Rand{r}
}

// Bytes generates n bytes
func (r *Rand) Bytes(n int) []byte {
	b := make([]byte, n)
	r.rand.Read(b) // its ok to ignore err https://golang.org/pkg/math/rand/#Read
	return b
}

// IPv4 returns a random IP address
func (r *Rand) IPv4() net.IP {
	return net.IP(r.Bytes(4))
}

// IntRange returns a non-negative pseudo-random number in [min,max)
func (r *Rand) IntRange(min, max int) int {
	return min + r.rand.Intn(max-min)
}

// Rune returns a rune from [start,stop] e.g. 'a', 'z'
func (r *Rand) Rune(start, stop rune) rune {
	return rune(r.IntRange(int(start), int(stop)+1))
}

// Alpha returns a random string of len(n) with azAZ chars.
func (r *Rand) Alpha(n int) string {
	return r.String(n, alpha)
}

// Letters returns a random string of n runes with alpha and a selection of utf8 letters.
func (r *Rand) Letters(n int) string {
	return r.String(n, letters)
}

func (r *Rand) String(n int, chars string) string {
	l := utf8.RuneCountInString(chars)
	runes := make([]rune, l)
	for i, r := range chars {
		runes[i] = r
	}

	var b bytes.Buffer
	for i := 0; i < l; i++ {
		b.WriteRune(runes[r.rand.Intn(l)])
	}
	return b.String()
}

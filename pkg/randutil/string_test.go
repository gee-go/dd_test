package randutil

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"unicode"
)

type unicodeTestCase struct {
	r    []rune
	size int
}

func (*unicodeTestCase) Generate(rand *rand.Rand, size int) reflect.Value {
	tc := &unicodeTestCase{
		r:    Quick(rand).Unicode(unicode.L, size),
		size: size,
	}

	return reflect.ValueOf(tc)
}

func TestUnicode(t *testing.T) {
	t.Parallel()
	rt := unicode.L

	// select every uint16 letter
	for i := 0; i < countR16(rt); i++ {
		r := selectR16(rt, i)
		if !unicode.IsLetter(r) {
			t.Fatalf("%c %v at %v should be a letter", r, r, i)
		}
	}

	// select every uint32 letter
	for i := 0; i < countR32(rt); i++ {
		r := selectR32(rt, i)
		if !unicode.IsLetter(r) {
			t.Fatalf("%c %U at %v should be a letter", r, r, i)
		}
	}

	// quick check unicode
	f := func(tc *unicodeTestCase) bool {
		if len(tc.r) != tc.size {
			return false
		}

		for _, r := range tc.r {
			if !unicode.IsLetter(r) {
				return false
			}
		}

		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

}

func BenchmarkUnicode(b *testing.B) {
	l := unicode.L
	r := New()

	for i := 0; i < b.N; i++ {
		r.Unicode(l, 1)
	}
}

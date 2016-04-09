package randutil

import (
	"testing"
	"unicode"
)

func TestString(t *testing.T) {
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

	// r := FromSeed(time.Now().UnixNano())
	// r.Unicode(unicode.L)
	// fmt.Printf("%c", r.Unicode(unicode.L))

	// // for _, r := range l.R16 {
	// //   for i
	// // }

}

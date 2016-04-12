package ddlog

import (
	"unicode"
	"unicode/utf8"
)

//buffer a utility to aid in parsing log lines
type buffer struct {
	s   string
	pos int
}

func (b *buffer) Init(s string) {
	b.s = s
	b.pos = 0
}

func (b *buffer) advanceUntil(end rune) string {
	var r rune
	var v string

	for i, w := b.pos, 0; i < len(b.s); i += w {
		r, w = utf8.DecodeRuneInString(b.s[i:])

		// special case for space to account for unicode.
		if r == end || (end == ' ' && unicode.IsSpace(r)) {
			v = b.s[b.pos:i]
			b.pos = i + w
			break
		}

	}

	if v == "" {
		v = b.s[b.pos:]
	}

	b.skipDelim()

	return v
}

func (b *buffer) skipDelim() {
	var r rune
	for i, w := b.pos, 0; i < len(b.s); i += w {
		r, w = utf8.DecodeRuneInString(b.s[i:])

		if isDelim(r) {
			b.pos += w
		} else {
			return
		}
	}
}

func isDelim(r rune) bool {
	switch r {
	case ' ', ']', '"', '[':
		return true
	}

	return unicode.IsSpace(r)
}

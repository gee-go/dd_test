package main

import (
	"unicode"

	"golang.org/x/exp/utf8string"
)

// Parser is used to convert lines to Message structs
type Parser struct {
	s   *utf8string.String
	f   string
	pos int
	end rune
}

func NewParser(f string) *Parser {
	return &Parser{f: f}
}

// parseUntil scans through the log line util the fn returns true.
func (p *Parser) parseUntil(fn func(rune) bool) string {
	var v string
	l := p.s.RuneCount()

	for i := p.pos; i < l; i++ {
		r := p.s.At(i)
		if fn(r) {
			v = p.s.Slice(p.pos, i)
			p.pos = i + 1
			break
		}
	}

	// at end of string
	if v == "" {
		return p.s.Slice(p.pos, l)
	}

	// skip delim chars.
	for i := p.pos; i < l; i++ {
		r := p.s.At(i)
		if unicode.IsSpace(r) || r == ']' || r == '"' || r == '[' {
			p.pos++
		} else {
			return v
		}
	}

	return v
}

func (p *Parser) parse() string {
	if unicode.IsSpace(p.end) {
		return p.parseUntil(unicode.IsSpace)
	}

	s := p.parseUntil(func(r rune) bool {
		return r == p.end
	})

	return s
}

// Parse converts a line to a message.
func (p *Parser) Parse(l string) *Message {
	p.pos = 0
	p.end = ' '
	p.s = utf8string.NewString(l)

	fieldStart := -1
	prev := ' '
	msg := &Message{}

	for i, r := range p.f {
		// field names
		switch r {
		case '{':
			switch prev {
			case '[':
				p.end = ']'
			case '"':
				p.end = '"'
			default:
				p.end = prev
			}
		case '}':
			msg.set(p.f[fieldStart:i], p.parse())
		}

		if prev == '{' {
			fieldStart = i
		}

		prev = r
	}

	return msg

}

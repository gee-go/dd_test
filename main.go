package main

import (
	"fmt"
	"unicode"

	"golang.org/x/exp/utf8string"
)

const (
	// DefaultLogFormat is a format string for the common log format
	DefaultLogFormat = `{remote} {ident} {auth} [{time}] "{request}" {status} {size}`

	// DefaultTimeFormat is the default format string used to parse timestamps
	DefaultTimeFormat = "02/Jan/2006:15:04:05 -0700"
)

type Parser struct {
	s   *utf8string.String
	pos int
}

func NewParser(s string) *Parser {
	return &Parser{s: utf8string.NewString(s)}
}

func (p *Parser) parseUntilEqual(end rune) string {
	s := p.parseUntil(func(r rune) bool {
		return r == end
	})

	p.pos++
	return s
}

func (p *Parser) parseUntil(fn func(rune) bool) string {
	var v string
	l := p.s.RuneCount()

	for i := p.pos; i < l; i++ {
		r := p.s.At(i)

		// consume until fn true or end of string
		if fn(r) || i == l-1 {
			v = p.s.Slice(p.pos, i)
			p.pos = i + 1
			break
		}
	}

	for i := p.pos; i < l; i++ {
		r := p.s.At(i)
		if unicode.IsSpace(r) {
			p.pos++
		} else {
			return v
		}
	}
	return v
}

func (p *Parser) parseToSpace() string {

	return p.parseUntil(unicode.IsSpace)

}

func (p *Parser) Parse(f string) {
	fieldStart := -1
	field := ""
	mode := ' '
	fmt.Println(p.s)
	for i, r := range f {
		switch {
		case r == '{':
			fieldStart = -1
		case r == '}':
			field = f[fieldStart:i]
		case r == '[':
			p.pos++
			mode = ']'
		case r == '"' && r != mode:
			p.pos++
			mode = '"'
		case unicode.IsLetter(r):
			if fieldStart == -1 {
				fieldStart = i
			}
		case r == mode:
			var v string
			switch mode {
			case ' ':
				v = p.parseToSpace()
			case ']', '"':
				v = p.parseUntilEqual(mode)
				mode = ' '
			}
			fmt.Println(field, v)

		default:
			fmt.Printf("%c %v\n", r, r)
		}
	}
}

// func ParseFormat(f string) {
// 	fieldStart := -1
// 	field := ""

//   for i, r := range f {
// 		switch {
// 		case r == '{':
// 			fieldStart = -1
// 		case r == '}':
// 			field = f[fieldStart:i]
// 		case unicode.IsLetter(r):
// 			if fieldStart == -1 {
// 				fieldStart = i
// 			}
// 		default:
// 			fmt.Printf("%c %v\n", r, r)
// 		}
// 	}
// }

// TODO - check for urls that contain [] or ""
func ParseLine(l string) *Message {
	na := 0
	fieldStart := -1
	mode := ' '

	msg := &Message{}

	for i, r := range l {
		switch r {
		case mode:
			if fieldStart >= 0 {
				msg.set(na, l[fieldStart:i])
				na++
				fieldStart = -1
			}
			mode = ' '
		case '[':
			if mode == ' ' {
				fieldStart = -1
				mode = ']'
			}

		case '"':
			if mode == ' ' {
				fieldStart = -1
				mode = '"'
			}
		default:
			if fieldStart == -1 {
				fieldStart = i
			}
		}
	}
	if fieldStart >= 0 { // Last field might end at EOF.
		msg.set(na, l[fieldStart:])
	}
	return msg
}

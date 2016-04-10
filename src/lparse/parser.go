package lparse

import (
	"unicode"
	"unicode/utf8"
)

type Line struct {
	Msg *Message
	Err error
}

// Parser is used to convert lines to Message structs
type Parser struct {
	s   string
	pos int
	end rune

	config *Config
}

func New(c *Config) *Parser {
	return &Parser{config: c}
}

// parseUntil scans through the log line util the fn returns true.
func (p *Parser) parseUntil(fn func(rune) bool) string {
	var v string
	var r rune

	// scan until fn return true
	for i, w := p.pos, 0; i < len(p.s); i += w {
		r, w = utf8.DecodeRuneInString(p.s[i:])

		if fn(r) {
			v = p.s[p.pos:i]
			p.pos = i + w
			break
		}

	}

	// at end of string
	if v == "" {
		return p.s[p.pos:]
	}

	// skip delim chars.
	for i, w := p.pos, 0; i < len(p.s); i += w {
		r, w = utf8.DecodeRuneInString(p.s[i:])

		if r == ' ' || r == ']' || r == '"' || r == '[' || unicode.IsSpace(r) {
			p.pos += w
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
func (p *Parser) Parse(l string) (*Message, error) {
	p.pos = 0
	p.end = ' '
	p.s = l

	fieldStart := -1
	prev := ' '
	msg := &Message{}

	for i, r := range p.config.LogFormat {
		// field names
		switch r {
		case '{':
			fieldStart = i + 1
			switch prev {
			case '[':
				p.end = ']'
			case '"':
				p.end = '"'
			default:
				p.end = prev
			}
		case '}':
			if err := msg.set(p.config.LogFormat[fieldStart:i], p.parse()); err != nil {
				return msg, err
			}
		}

		prev = r
	}

	return msg, nil

}

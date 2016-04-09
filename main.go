package main

import (
	"fmt"
	"unicode"
)

const (
	// DefaultLogFormat is a format string for the common log format
	DefaultLogFormat = `{remote} {ident} [{time}] "{method} {uri} {proto}" {status} {size}`

	// DefaultTimeFormat is the default format string used to parse timestamps
	DefaultTimeFormat = "02/Jan/2006:15:04:05 -0700"
)

func ParseFormat(f string) {
	fieldStart := -1

	for i, r := range f {
		switch {
		case r == '{':
			fieldStart = -1
		case r == '}':
			fmt.Println(f[fieldStart:i])
		case unicode.IsLetter(r):
			if fieldStart == -1 {
				fieldStart = i
			}
		default:
			fmt.Printf("%c %v\n", r, r)
		}
	}
}

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

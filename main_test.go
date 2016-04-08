package main

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

const example = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

var delimMap = map[rune]rune{
	'[': ']',
	'"': '"',
}

type Field struct {
	Name    string
	EndRune rune
}

func LogFields2(l, f string) map[string]string {
	out := make(map[string]string)

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

	return out
}

func BenchmarkFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Fields(example)
	}
}

func TestLogFields2(t *testing.T) {
	f := `{remote} {ident} [{time}] "{method} {uri} {proto}" {status} {size}`

	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

	LogFields2(l, f)

}

func BenchmarkLogFields(b *testing.B) {
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	for i := 0; i < b.N; i++ {
		ParseLine(l)
	}
}

func TestMain(t *testing.T) {
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	pp.Println(ParseLine(l))

	l2 := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif?a=[ HTTP/1.0" 200 2326`
	msg := ParseLine(l2)
	assert := require.New(t)
	assert.Equal("GET /apache_pb.gif?a=[ HTTP/1.0", msg.Method)

}

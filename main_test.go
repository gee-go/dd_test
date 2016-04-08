package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

type Parser struct {
	timeFormat string
}

type Message struct {
	Remote        string
	Ident         string
	Auth          string
	Time          time.Time
	Request       string
	Status        string
	BodyBytesSent string
}

func (m *Message) set(i int, s string) {
	switch i {
	case 3:
		t, err := time.Parse("02/Jan/2006:15:04:05 -0700", s)
		m.Time = t
		if err != nil {
			panic(err)
		}
	default:
		reflect.ValueOf(m).Elem().Field(i).SetString(s)
	}
}

var delimMap = map[rune]rune{
	'[': ']',
	'"': '"',
}

// TODO - check for urls that contain [] or ""
func LogFields(l string) *Message {
	// a := make([]string, 7)
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

	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	for i := 0; i < b.N; i++ {
		strings.Fields(l)
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
		LogFields(l)
	}
}

func TestMain(t *testing.T) {
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	pp.Println(LogFields(l))

	l2 := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif?a=[ HTTP/1.0" 200 2326`
	msg := LogFields(l2)
	assert := require.New(t)
	assert.Equal("GET /apache_pb.gif?a=[ HTTP/1.0", msg.Request)

}

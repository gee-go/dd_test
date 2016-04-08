package main

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

type Parser struct {
	timeFormat string
}

type Message struct {
	RemoteAddr    string
	UserId        string
	RemoteUser    string
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

func LogFields(l string) *Message {
	// a := make([]string, 7)
	na := 0
	fieldStart := -1
	mode := ' '

	msg := &Message{}
	// rmsg := reflect.ValueOf(msg).Elem()

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
			fieldStart = -1
			mode = ']'
		case '"':
			fieldStart = -1
			mode = '"'
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

func LogFields2(l, f string) *Message {
	msg := &Message{}

	return msg
}

func BenchmarkFields(b *testing.B) {

	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	for i := 0; i < b.N; i++ {
		strings.Fields(l)
	}
}

func BenchmarkLogFields2(b *testing.B) {
	f := `%h %l %u [%t] "%r" %>s %b`
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	for i := 0; i < b.N; i++ {
		LogFields2(l, f)
	}
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
}

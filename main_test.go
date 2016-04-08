package main

import (
	"strings"
	"testing"

	"github.com/k0kubun/pp"
)

type Message struct {
	RemoteAddr string
	UserId     string
	UserName   string
	Time       string
	Request    string
	StatusCode string
	Size       string
}

func LogFields(l string) []string {
	a := make([]string, 7)
	na := 0
	fieldStart := -1
	mode := ' '

	for i, r := range l {
		switch r {
		case mode:
			if fieldStart >= 0 {
				a[na] = l[fieldStart:i]
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
		a[na] = l[fieldStart:]
	}
	return a
}

func BenchmarkFields(b *testing.B) {
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	for i := 0; i < b.N; i++ {
		strings.Fields(l)
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

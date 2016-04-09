package main

import "testing"

const example = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

func BenchmarkParser(b *testing.B) {
	p := NewParser(DefaultLogFormat)

	for i := 0; i < b.N; i++ {
		p.Parse(example)
	}
}

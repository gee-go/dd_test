package main

import "testing"

const example = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

func TestMain(t *testing.T) {
	l := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`
	p := NewParser(l)
	p.Parse(DefaultLogFormat)

	// l2 := `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif?a=[ HTTP/1.0" 200 2326`
	// msg := ParseLine(l2)
	// assert := require.New(t)
	// assert.Equal("GET /apache_pb.gif?a=[ HTTP/1.0", msg.Request)

}

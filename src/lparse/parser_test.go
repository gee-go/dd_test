package lparse

import (
	"fmt"
	"testing"
	"time"
)

func newParser(f string) *Parser {
	c := NewConfig()
	c.LogFormat = f
	return New(c)
}

func TestParserExample(t *testing.T) {
	t.Parallel()
	tc, a := newTestCase(t)

	m, err := tc.p.Parse(ExampleLogLine)
	a.NoError(err)
	a.Equal(m.Remote, "127.0.0.1")
	a.Equal(m.Ident, "-")
	a.Equal(m.Auth, "frank")
	a.Equal(m.Time.Year(), 2000)
	a.Equal(m.Method, "GET")
	a.Equal(m.URI, "/apache_pb.gif")
	a.Equal(m.Proto, "HTTP/1.0")
	a.Equal(m.Status, "200")
	a.Equal(m.Size, "2326")
}

func TestParserFormat(t *testing.T) {
	t.Parallel()
	// swapped ident and remote
	tc, a := newTestCase(t, `{ident} {remote} {auth} [{time}] "{request}" {status} {size}`)

	m := tc.MustParse(ExampleLogLine)
	a.Equal(m.Ident, "127.0.0.1")
	a.Equal(m.Remote, "-")
}

// TODO
// func TestParserFormatStrange(t *testing.T) {
// 	// starts with time
// 	t.Parallel()
// 	tc, _ := newTestCase(t, `[{time}]`)
// 	v := time.Now()

// 	tc.MustParse(fmt.Sprintf("[%s]", v.Format(tc.c.TimeFormat)))
// }

func TestParserError(t *testing.T) {
	t.Parallel()
	tc, a := newTestCase(t, `{ident} [{time}]`)

	// make sure it works.
	tc.MustParse(fmt.Sprintf("- [%s]", time.Now().Format(tc.c.TimeFormat)))

	_, err := tc.p.Parse("- [abc]")
	a.Error(err)
	a.Equal("abc", err.(*time.ParseError).Value)
}

func TestRandomMessages(t *testing.T) {
	t.Parallel()
	tc, a := newTestCase(t)

	for i := 0; i < 1000; i++ {
		m := tc.g.RandMsg()
		pm := tc.MustParse(m.String())

		a.Equal(m, pm)

	}
}

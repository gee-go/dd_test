package ddlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// common test setup
type testCase struct {
	p *Parser
	g *Generator
	c *Config
	t *testing.T
	a *require.Assertions
}

func newTestCase(t *testing.T, fmt ...string) (*testCase, *require.Assertions) {
	c := NewConfig()

	if len(fmt) == 1 {
		c.LogFormat = fmt[0]
	}
	assert := require.New(t)
	return &testCase{
		p: c.NewParser(),
		g: c.NewGenerator(),
		c: c,
		t: t,
		a: assert,
	}, assert
}

func (tc *testCase) MustParse(l string) *Message {
	m, err := tc.p.Parse(l)
	tc.a.NoError(err)
	return m
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
	a.Equal(m.Status, 200)
	a.Equal(m.Size, 2326)
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
	m := tc.g.TestMsg()
	line := m.Format(tc.c.LogFormat)
	tc.MustParse(line)

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

package lparse

import (
	"reflect"
	"testing"

	"github.com/gee-go/dd_test/src/randutil"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/require"
)

const exampleLine = `127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326`

func newParser(f string) *Parser {
	c := NewConfig()
	c.LogFormat = f
	return New(c)
}

func TestParserExample(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	p := newParser(DefaultLogFormat)

	m, err := p.Parse(exampleLine)
	assert.NoError(err)
	assert.Equal(m.Remote, "127.0.0.1")
	assert.Equal(m.Ident, "-")
	assert.Equal(m.Auth, "frank")
	assert.Equal(m.Time.Year(), 2000)
	assert.Equal(m.Method, "GET")
	assert.Equal(m.URI, "/apache_pb.gif")
	assert.Equal(m.Proto, "HTTP/1.0")
	assert.Equal(m.Status, "200")
	assert.Equal(m.Size, "2326")
}

func TestParserFormat(t *testing.T) {
	t.Parallel()
	// swapped ident and remote
	assert := require.New(t)
	p := newParser(`{ident} {remote} {auth} [{time}] "{request}" {status} {size}`)
	m, err := p.Parse(exampleLine)
	assert.NoError(err)
	assert.Equal(m.Ident, "127.0.0.1")
	assert.Equal(m.Remote, "-")
}

func BenchmarkParser(b *testing.B) {
	p := newParser(DefaultLogFormat)

	for i := 0; i < b.N; i++ {
		p.Parse(exampleLine)
	}
}

func TestRandomMessages(t *testing.T) {
	t.Parallel()
	p := newParser(DefaultLogFormat)

	for i := 0; i < 1000; i++ {
		m := randMessage(randutil.R)
		pm, err := p.Parse(m.String())
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(m, pm) {
			pp.Println(m, pm)
			t.Fatalf("%v %v", m, pm)
		}

	}
}

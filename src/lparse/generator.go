package lparse

import (
	"bytes"
	"math/rand"
	"net/url"
	"strconv"
	"time"
	"unicode"

	"github.com/gee-go/dd_test/src/randutil"
)

// A Generator is used to create random messages.
type Generator struct {
	c *Config
	r *randutil.Rand

	// If true, Generate fields with full range of unicode chars.
	// Otherwise only use ascii.
	UseUnicode bool
}

// NewGenerator inits a Generator with the given config.
// If config is nil, it uses the default config.
func NewGenerator(c *Config) *Generator {
	if c == nil {
		c = NewConfig()
	}

	return &Generator{
		c:          c,
		r:          randutil.R,
		UseUnicode: true,
	}
}

func (g *Generator) randURI() string {
	r := g.r
	var path bytes.Buffer

	// 1 to 5 path components
	for i := 0; i < r.IntRange(1, 5); i++ {
		path.WriteString("/")
		path.WriteString(r.Alpha(r.Rand.Intn(6)))
	}

	u, _ := url.Parse(path.String())

	// 0 - 5 random url params
	for i := 0; i < r.Rand.Intn(6); i++ {
		u.Query().Set(r.Alpha(r.IntRange(1, 5)), r.Alpha(r.IntRange(1, 5)))
	}

	return u.RequestURI()
}

func (g *Generator) randAuth(maxLen int) string {
	var auth string
	authLen := rand.Intn(maxLen)
	if g.UseUnicode {
		auth = string(g.r.Unicode(unicode.L, authLen))
	} else {
		auth = g.r.Alpha(authLen)
	}

	if len(auth) == 0 {
		auth = "-"
	}

	return auth
}

// RandMsg creates a random Message
func (g *Generator) RandMsg() *Message {
	t, _ := time.Parse(g.c.TimeFormat, time.Now().Format(g.c.TimeFormat))

	return &Message{
		Remote: g.r.IPv4().String(),
		Ident:  "-",
		Auth:   g.randAuth(20),
		Time:   t,
		Method: g.r.SelectString("GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"),
		URI:    g.randURI(),
		Proto:  "HTTP/1.0",
		Status: g.r.SelectString("200", "400", "201", "304", "401", "404", "500"),
		Size:   strconv.Itoa(g.r.IntRange(1<<8, 1<<26)),
	}
}

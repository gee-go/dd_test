package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gee-go/dd_test/pkg/randutil"
)

type Request struct {
	Method string
	URI    string
	Proto  string
}

func randRequest(r *randutil.Rand) *Request {
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

	return &Request{
		Method: r.SelectString("GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"),
		URI:    u.String(),
		Proto:  "HTTP/1.0",
	}

}

func (r *Request) String() string {
	return fmt.Sprintf(`"%s %s %s"`, r.Method, r.URI, r.Proto)
}

// Message represents a single log line.
type Message struct {
	Remote  string    `log:"remote"`
	Ident   string    `log:"ident"`
	Auth    string    `log:"auth"`
	Time    time.Time `log:"time"`
	Request *Request  `log:"request"`
	Status  string    `log:"status"`
	Size    string    `log:"size"`
}

func (m *Message) String() string {
	t := fmt.Sprintf("[%s]", m.Time.Format(DefaultTimeFormat))
	parts := []string{m.Remote, m.Ident, m.Auth, t, m.Request.String(), m.Status, m.Size}
	return strings.Join(parts, " ")
}

func (m *Message) set(f, s string) error {
	var err error
	switch f {
	case "remote":
		m.Remote = s
	case "ident":
		m.Ident = s
	case "auth":
		m.Auth = s
	case "time":
		m.Time, err = time.Parse(DefaultTimeFormat, s)
	case "request":
		parts := strings.Fields(s)
		if len(parts) != 3 {
			return fmt.Errorf("%s is an invalid request", s)
		}

		m.Request = &Request{
			Method: parts[0],
			URI:    parts[1],
			Proto:  parts[2],
		}
	case "status":
		m.Status = s
	case "size":
		m.Size = s
	default:
		return fmt.Errorf("Unknown field %s", f)
	}

	return err
}

func randMessage(r *randutil.Rand) *Message {
	auth := string(r.Unicode(unicode.L, rand.Intn(20)))

	if len(auth) == 0 {
		auth = "-"
	}

	// format then parse so it has the same precision as the parser
	t, _ := time.Parse(DefaultTimeFormat, time.Now().Format(DefaultTimeFormat))

	return &Message{
		Remote:  r.IPv4().String(),
		Ident:   "-",
		Auth:    auth,
		Time:    t,
		Request: randRequest(r),
		Status:  r.SelectString("200", "400"),
		Size:    strconv.Itoa(r.IntRange(1<<8, 1<<26)),
	}
}

// Generate a random message for using the testing/quick package
func (*Message) Generate(rand *rand.Rand, size int) reflect.Value {
	r := randutil.Quick(rand)

	return reflect.ValueOf(randMessage(r))
}

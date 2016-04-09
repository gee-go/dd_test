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

func (r *Request) Format() string {
	return fmt.Sprintf("%s %s %s", r.Method, url.QueryEscape(r.URI), r.Proto)
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

func (m *Message) Format() string {
	v := reflect.ValueOf(m).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Println(field.Tag.Get("log"))
	}

	return ""
}

func randRequest(r *randutil.Rand) *Request {
	var path bytes.Buffer
	// path components
	for i := 0; i < r.IntRange(1, 5); i++ {
		path.WriteString(r.Alpha(r.IntRange(1, 5)))
	}

	u, _ := url.Parse(path.String())

	return &Request{
		Method: r.SelectString("GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"),
		URI:    u.RequestURI(),
		Proto:  "HTTP/1.0",
	}

}

// Generate a random message for using the testing/quick package
func (*Message) Generate(rand *rand.Rand, size int) reflect.Value {
	r := randutil.Quick(rand)

	auth := string(r.Unicode(unicode.L, rand.Intn(20)))
	if len(auth) == 0 {
		auth = "-"
	}

	// format then parse so it has the same precision as the parser
	t, _ := time.Parse(DefaultLogFormat, time.Now().Format(DefaultLogFormat))

	m := &Message{
		Remote:  r.IPv4().String(),
		Ident:   "-",
		Auth:    auth,
		Time:    t,
		Request: randRequest(r),
		Status:  r.SelectString("200", "400"),
		Size:    strconv.Itoa(r.IntRange(1<<8, 1<<26)),
	}

	return reflect.ValueOf(m)
}

func (m *Message) set(i int, s string) error {
	var err error
	switch i {
	case 3:
		m.Time, err = time.Parse(DefaultTimeFormat, s)
	case 4:
		parts := strings.Fields(s)
		if len(parts) != 3 {
			return fmt.Errorf("%s is an invalid request", s)
		}

		m.Request = &Request{
			Method: parts[0],
			URI:    parts[1],
			Proto:  parts[2],
		}
	default:
		reflect.ValueOf(m).Elem().Field(i).SetString(s)
	}

	return err
}

package main

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/gee-go/dd_test/pkg/randutil"
)

// Message represents a single log line.
type Message struct {
	Remote string    `log:"remote"`
	Ident  string    `log:"ident"`
	Auth   string    `log:"auth"`
	Time   time.Time `log:"time"`
	Method string    `log:"method"`
	URI    string    `log:"uri"`
	Proto  string    `log:"proto"`
	Status string    `log:"status"`
	Size   string    `log:"size"`
}

// Generate a random message for using the testing/quick package
func (*Message) Generate(rand *rand.Rand, size int) reflect.Value {
	m := &Message{}
	r := randutil.New(rand)
	m.Remote = r.IPv4().String()
	m.Ident = "-"
	return reflect.ValueOf(m)
}

func (m *Message) set(i int, s string) {
	switch i {
	case 3:
		t, err := time.Parse(DefaultTimeFormat, s)
		m.Time = t
		if err != nil {
			panic(err)
		}
	default:
		reflect.ValueOf(m).Elem().Field(i).SetString(s)
	}
}

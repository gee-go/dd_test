package main

import (
	"reflect"
	"testing"
	"testing/quick"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	p := NewParser(DefaultLogFormat)

	f := func(m *Message) bool {
		l := m.Format()
		return reflect.DeepEqual(m, p.Parse(l))
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

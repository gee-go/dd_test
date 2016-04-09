package main

import (
	"testing"
	"testing/quick"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	f := func(m *Message) bool {
		m.Format()
		return false
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

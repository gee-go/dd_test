package main

import (
	"testing"

	"github.com/gee-go/dd_test/pkg/randutil"
	"github.com/stretchr/testify/require"
)

func TestMessage(t *testing.T) {
	t.Parallel()
	p := NewParser(DefaultLogFormat)
	assert := require.New(t)

	for i := 0; i < 1000; i++ {
		m := randMessage(randutil.R)
		assert.Equal(m, p.Parse(m.String()))
	}
}

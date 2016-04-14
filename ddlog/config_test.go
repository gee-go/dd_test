package ddlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	a := require.New(t)

	mc := &Config{
		WindowSize:  2 * time.Minute,
		AggInterval: 1 * time.Second,
	}

	a.Equal(120, mc.numWindowsKept())
}

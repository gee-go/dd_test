package ddlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAlert(t *testing.T) {
	alert := &Alert{
		Start: time.Now(),
		Count: 5,
	}
	a := require.New(t)
	a.False(alert.IsDone())

	alert.Complete(time.Now().Add(time.Second))
	a.True(alert.IsDone())
}

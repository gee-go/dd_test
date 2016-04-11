package lparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageEventName(t *testing.T) {
	cases := []struct {
		uri  string
		name string
	}{
		{"/", "/"},
		{"/a", "/a"},
		{"/a/a", "/a"},
		{"/ab/a", "/ab"},
		{"/ab?a=1", "/ab"},
		{"http://my.site.com/pages/create", "http://my.site.com/pages"},
	}

	a := require.New(t)
	for _, tc := range cases {
		m := &Message{URI: tc.uri}
		a.Equal(tc.name, m.EventName())
	}
}

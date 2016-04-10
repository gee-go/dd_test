package lparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// common test setup
type testCase struct {
	p *Parser
	g *Generator
	c *Config
	t *testing.T
	a *require.Assertions
}

func newTestCase(t *testing.T, fmt ...string) (*testCase, *require.Assertions) {
	c := NewConfig()
	if len(fmt) == 1 {
		c.LogFormat = fmt[0]
	}
	assert := require.New(t)
	return &testCase{
		p: New(c),
		g: NewGenerator(c),
		c: c,
		t: t,
		a: assert,
	}, assert
}

func (tc *testCase) MustParse(l string) *Message {
	m, err := tc.p.Parse(l)
	tc.a.NoError(err)
	return m
}

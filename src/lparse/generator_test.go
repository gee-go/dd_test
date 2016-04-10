package lparse

import "testing"

func TestGeneratorTestMsg(t *testing.T) {
	t.Parallel()
	tc, a := newTestCase(t)

	a.Equal(ExampleLogLine, tc.g.TestMsg().String())
}

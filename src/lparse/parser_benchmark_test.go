package lparse

import "testing"

func benchLineParser(b *testing.B, fmt, line string) {
	p := newParser(fmt)

	// error check
	_, err := p.Parse(line)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Parse(line)
	}
}

func BenchmarkParserSimple(b *testing.B) {
	benchLineParser(b, DefaultLogFormat, ExampleLogLine)
}

func BenchmarkParserUnicode(b *testing.B) {
	benchLineParser(b, DefaultLogFormat, ExampleLogLine)
}

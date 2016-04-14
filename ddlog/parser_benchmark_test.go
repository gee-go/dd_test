package ddlog

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func benchLineParser(b *testing.B, fmt, line string) {
	p := NewConfig().NewParser()

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
	b.ReportAllocs()
	benchLineParser(b, DefaultLogFormat, ExampleLogLine)
}

func BenchmarkParserUnicode(b *testing.B) {
	benchLineParser(b, DefaultLogFormat, ExampleLogLineUnicode)
}

// regex version for comparison.
func formatToReg(format string) *regexp.Regexp {
	var parts []string
	fieldStart := -1

	for i, r := range format {
		switch r {
		case '{':
			fieldStart = i + 1
		case '}':
			field := format[fieldStart:i]
			switch field {
			case "time":
				parts = append(parts, `\[(?P<time>[^]]+)\]`)
			case "request":
				parts = append(parts, `"(?P<request>[^"]+)"`)
			default:
				parts = append(parts, fmt.Sprintf(`(?P<%s>[^ ]+)`, field))
			}
		}

	}
	return regexp.MustCompile(fmt.Sprintf("^%s$", strings.Join(parts, " ")))
}

func parseRegex(re *regexp.Regexp, l string) *Message {
	fields := re.FindStringSubmatch(l)
	msg := &Message{}
	for i, name := range re.SubexpNames() {
		if i == 0 {
			continue
		}
		msg.set(name, fields[i])

	}

	return msg
}

func BenchmarkParserRegex(b *testing.B) {
	re := formatToReg(DefaultLogFormat)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parseRegex(re, ExampleLogLine)
	}
}

func BenchmarkParserRegexUnicode(b *testing.B) {
	re := formatToReg(DefaultLogFormat)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		parseRegex(re, ExampleLogLineUnicode)
	}
}

package lparse

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// Message represents a single log line.
type Message struct {
	Remote string
	Ident  string
	Auth   string
	Time   time.Time
	Method string
	URI    string
	Proto  string
	Status string
	Size   string
}

func (m *Message) get(f string) string {
	switch f {
	case "remote":
		return m.Remote
	case "ident":
		return m.Ident
	case "auth":
		return m.Auth
	case "time":
		return m.Time.Format(DefaultTimeFormat)
	case "request":
		return strings.Join([]string{m.Method, m.URI, m.Proto}, " ")
	case "status":
		return m.Status
	case "size":
		return m.Size
	default:

	}

	return "<err>"
}

// Format writes the message in the given format
func (m *Message) Format(format string) string {
	var b bytes.Buffer
	fieldStart := -1
	for i, r := range format {
		switch r {
		case '{':
			fieldStart = i + 1
		case '}':
			b.WriteString(m.get(format[fieldStart:i]))
			fieldStart = -1
		default:
			if fieldStart == -1 {
				b.WriteRune(r)
			}
		}
	}

	return b.String()
}

func (m *Message) String() string {
	return m.Format(DefaultLogFormat)
}

func (m *Message) set(f, s string) error {
	var err error
	switch f {
	case "remote":
		m.Remote = s
	case "ident":
		m.Ident = s
	case "auth":
		m.Auth = s
	case "time":
		m.Time, err = time.Parse(DefaultTimeFormat, s)
	case "request":
		parts := strings.Fields(s)
		if len(parts) != 3 {
			return fmt.Errorf("%s is an invalid request", s)
		}

		m.Method = parts[0]
		m.URI = parts[1]
		m.Proto = parts[2]

	case "status":
		m.Status = s
	case "size":
		m.Size = s
	default:
		return fmt.Errorf("Unknown field %s", f)
	}

	return err
}

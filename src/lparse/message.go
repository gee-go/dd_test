package lparse

import (
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

func (m *Message) String() string {
	t := fmt.Sprintf("[%s]", m.Time.Format(DefaultTimeFormat))
	req := fmt.Sprintf(`"%s %s %s"`, m.Method, m.URI, m.Proto)
	parts := []string{m.Remote, m.Ident, m.Auth, t, req, m.Status, m.Size}
	return strings.Join(parts, " ")
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

package ddlog

import (
	"time"

	"github.com/benbjohnson/clock"
)

const (
	// DefaultLogFormat is a format string for the common log format
	DefaultLogFormat = `{remote} {ident} {auth} [{time}] "{request}" {status} {size}`

	// DefaultTimeFormat is the default format string used to parse timestamps
	DefaultTimeFormat = "02/Jan/2006:15:04:05 -0700"
)

type Config struct {
	clock clock.Clock

	// Parse Config
	LogFormat  string
	TimeFormat string

	// Montitor Config
	AggInterval    time.Duration // How often to aggregate lines for the rolling window.
	WindowSize     time.Duration // How long to keep aggregates for the rolling window.
	AlertThreshold int           // If number of messages over the past WindowSize exceeds this number trigger and alert.

	// File Config
	Filename string
}

// Mock time for all things.
func (c *Config) Mock(cl clock.Clock) {
	c.clock = cl
}

func NewConfig() *Config {
	return &Config{
		clock:          clock.New(),
		LogFormat:      DefaultLogFormat,
		TimeFormat:     DefaultTimeFormat,
		AggInterval:    time.Second * 1,
		WindowSize:     time.Minute * 2,
		AlertThreshold: 100,
	}
}

func (c *Config) numWindowsKept() int {
	return int(c.WindowSize / c.AggInterval)
}

func (c *Config) NewMonitor() *Monitor {
	return NewMonitor(c)
}

func (c *Config) NewParser() *Parser {
	return NewParser(c)
}

func (c *Config) NewGenerator() *Generator {
	return NewGenerator(c)
}

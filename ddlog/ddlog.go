package ddlog

import "time"

const (
	// DefaultLogFormat is a format string for the common log format
	DefaultLogFormat = `{remote} {ident} {auth} [{time}] "{request}" {status} {size}`

	// DefaultTimeFormat is the default format string used to parse timestamps
	DefaultTimeFormat = "02/Jan/2006:15:04:05 -0700"
)

type Config struct {
	LogFormat        string
	TimeFormat       string
	Filename         string
	FastTickDuration time.Duration
	AlertDuration    time.Duration
	AlertThreshold   int
}

func NewConfig() *Config {
	return &Config{
		LogFormat:        DefaultLogFormat,
		TimeFormat:       DefaultTimeFormat,
		FastTickDuration: time.Second * 10,
		AlertDuration:    time.Minute * 2,
		AlertThreshold:   100,
	}
}

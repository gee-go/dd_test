package main

import (
	"flag"
	"time"

	"github.com/gee-go/dd_test/src/lparse"
)

type Opts struct {
	lparse.Config

	// How often to send logs.
	Rate time.Duration
}

func parseFlags() *Opts {
	o := &Opts{}
	flag.StringVar(&o.LogFormat, "fmt", lparse.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", lparse.DefaultTimeFormat, "a")

	return o
}

func main() {

}

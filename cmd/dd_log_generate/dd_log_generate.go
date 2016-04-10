package main

import (
	"flag"
	"fmt"
	"log"
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
	var rate string
	flag.StringVar(&rate, "rate", "1s", "a")
	flag.Parse()

	r, err := time.ParseDuration(rate)
	if err != nil {
		log.Fatal(err)
	}
	o.Rate = r

	return o
}

func main() {
	opts := parseFlags()
	g := lparse.NewGenerator(&opts.Config)
	g.UseUnicode = false

	for range time.Tick(opts.Rate) {
		fmt.Println(g.RandMsg())
	}
}

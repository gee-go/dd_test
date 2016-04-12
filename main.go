package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/gee-go/dd_test/src/lscan"
	"github.com/gee-go/dd_test/src/metric"
)

func parseFlags() *lparse.Config {
	o := &lparse.Config{}
	flag.StringVar(&o.LogFormat, "fmt", lparse.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", lparse.DefaultTimeFormat, "a")

	flag.Parse()

	return o
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	config := parseFlags()

	s, err := lscan.Tail(flag.Args()[0], config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-c
		s.Cleanup()
		os.Exit(1)
	}()

	go s.Start()

	g := metric.NewGroup()
	go g.Start(s.MsgChan)

	done := make(chan bool)

	// window := lscan.NewWindow(3)

	// windowSize := time.Second * 10

	<-done
}

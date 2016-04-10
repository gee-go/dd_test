package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/k0kubun/pp"
)

func parseFlags() *lparse.Config {
	o := lparse.NewConfig()
	flag.StringVar(&o.LogFormat, "fmt", lparse.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", lparse.DefaultTimeFormat, "a")

	flag.Parse()

	return o
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	s := lparse.NewFileScanner(parseFlags())
	go func() {
		<-c
		s.Cleanup()
		os.Exit(1)
	}()

	for _, fn := range flag.Args() {
		go s.Tail(fn)
	}

	for l := range s.Line() {
		pp.Println(l)
	}

}

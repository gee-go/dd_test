package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/gee-go/dd_test/src/lscan"
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

	s := lscan.NewFileScanner(parseFlags())
	go func() {
		<-c
		s.Cleanup()
		os.Exit(1)
	}()

	for _, fn := range flag.Args() {
		go s.Tail(fn)
	}

	ms := lscan.NewMetric()

	for l := range s.Line() {
		if l.Err != nil {
			fmt.Println(l.Err)
			continue
		}

		ms.HandleMsg(l.Msg)
	}

}

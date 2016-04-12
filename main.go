package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VividCortex/ewma"
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

	// window := lscan.NewWindow(3)

	// windowSize := time.Second * 10

	tickChan := time.Tick(time.Second * 1)
	count := 0
	avg := ewma.NewMovingAverage()
	for {
		select {
		case <-tickChan:

			avg.Add(float64(count))
			fmt.Println(avg.Value())
			count = 0
		case <-s.MsgChan:
			// fmt.Println(m.Remote)
			count++
		}
	}

	// <-done
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

type MetricGroup struct {
	count int
	wTime time.Duration
	mu    sync.Mutex
}

func NewMetricGroup(wTime time.Duration) *MetricGroup {
	return &MetricGroup{wTime: wTime}
}

func start(s *lscan.TailScanner) {
	tickChan := time.Tick(time.Second * 5)

	ms := metric.New()

	for {
		select {
		case <-tickChan:
			for _, p := range ms.TopK(10) {
				fmt.Println(p.ID, p.Count)
			}

		case m := <-s.MsgChan:
			ms.HandleMsg(m)
		}
	}

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
	go start(s)

	done := make(chan bool)

	// window := lscan.NewWindow(3)

	// windowSize := time.Second * 10

	<-done
}

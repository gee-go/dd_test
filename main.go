package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gee-go/dd_test/ddlog"
	"github.com/gee-go/dd_test/src/lscan"
	"github.com/k0kubun/pp"
)

func parseFlags() *ddlog.Config {
	o := &ddlog.Config{}
	flag.StringVar(&o.LogFormat, "fmt", ddlog.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", ddlog.DefaultTimeFormat, "a")

	flag.Parse()

	return o
}

type Bucket struct {
	Time time.Time

	Count       int
	CountByPage map[string]int
}

func (b *Bucket) Add(m *ddlog.Message) {
	b.Count++
	name := m.EventName()
	if _, found := b.CountByPage[name]; !found {
		b.CountByPage[name] = 1
	} else {
		b.CountByPage[name]++
	}
}

func newBucket(mt time.Time) *Bucket {
	return &Bucket{
		CountByPage: make(map[string]int),
		Time:        mt,
	}
}

func BucketMsgChan(msgChan <-chan *ddlog.Message, bw time.Duration) chan *Bucket {
	out := make(chan *Bucket)

	go func() {
		var b *Bucket

		for m := range msgChan {
			mt := m.Time.Truncate(bw)

			if b == nil {
				// first bucket
				b = newBucket(mt)
			} else if mt.After(b.Time) {
				// new bucket
				out <- b
				b = newBucket(mt)
			}

			if mt.Equal(b.Time) {
				// same bucket
				b.Add(m)
			} else {
				// TODO - Handle case where log timestamps are not constantly increasing.
				fmt.Println("Old line")
			}

		}
		close(out)
	}()

	return out
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

	bChan := BucketMsgChan(s.MsgChan, time.Second*1)
	for b := range bChan {
		pp.Println(b)
	}

}

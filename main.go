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
	"github.com/hpcloud/tail"
)

func parseFlags() *ddlog.Config {
	o := &ddlog.Config{}
	flag.StringVar(&o.LogFormat, "fmt", ddlog.DefaultLogFormat, "a")
	flag.StringVar(&o.TimeFormat, "time", ddlog.DefaultTimeFormat, "a")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Fatal("Need a file")
	}

	o.Filename = args[0]

	return o
}

func main() {
	config := parseFlags()

	// setup file tailer
	fileTail, err := tail.TailFile(config.Filename, tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// setup tail cleanup
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		fileTail.Cleanup()
		os.Exit(0)
	}()

	// tail lines -> messages
	msgChan := make(chan *ddlog.Message)
	lineParser := ddlog.New(config)
	go func() {
		defer close(msgChan)
		for line := range fileTail.Lines {
			if line.Err != nil {
				fmt.Println(line.Err)
				continue
			}
			m, err := lineParser.Parse(line.Text)
			if err != nil {
				fmt.Println(err)
				continue
			}
			msgChan <- m
		}
	}()

	// messages -> metrics
	metricStore := ddlog.NewMetricStore()
	go metricStore.Start(msgChan)

	tick10s := time.Tick(time.Second * 10)
	tick2m := time.Tick(time.Minute * 2)

	for {
		select {
		case <-tick10s:
			metricStore.Print()
		case <-tick2m:
			fmt.Println("2min")
		}
	}
}

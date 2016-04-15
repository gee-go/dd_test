package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/gee-go/ddlog/ddlog"
	"github.com/gee-go/ddlog/ddlog/cli"
	"github.com/hpcloud/tail"
)

func parseFlags() *ddlog.Config {
	o := ddlog.NewConfig()
	flag.StringVar(&o.LogFormat, "fmt", ddlog.DefaultLogFormat, "Log format to parse")
	flag.StringVar(&o.TimeFormat, "time", ddlog.DefaultTimeFormat, "Time format to parse")
	flag.DurationVar(&o.WindowSize, "window", o.WindowSize, "Duration to monitor alert count over")
	flag.IntVar(&o.AlertThreshold, "alert", o.AlertThreshold, "Trigger alert when visit count exceeds this number over the given window")
	flag.BoolVar(&o.PlainUI, "plain", o.PlainUI, "Use non-fancy output")
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
	defer fileTail.Cleanup()

	// setup tail cleanup
	// quitChan := make(chan bool, 1)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		os.Exit(0)
	}()

	// tail lines -> messages
	msgChan := make(chan *ddlog.Message)
	lineParser := config.NewParser()
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

	mon := config.NewMonitor()
	go mon.Start(context.Background(), msgChan)

	if config.PlainUI {
		alertTicker := time.Tick(time.Second)
		topKTicker := time.Tick(5 * time.Second)
		table := cli.NewTable()
		// TODO - its possible a cycle might be skipped.
		for {
			select {
			case <-topKTicker:
				topk := mon.TopK(20)
				if len(topk) == 0 {
					continue
				}
				cli.SetTopK(table, topk)
				fmt.Println(table.String())
			case <-alertTicker:
				alert := mon.PopAlert()
				if alert != nil {
					fmt.Println(alert)
				}

			}
		}

	} else {
		ui := cli.NewUI(mon)
		ui.Start()
	}

}

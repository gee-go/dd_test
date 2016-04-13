package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gee-go/dd_test/ddlog"
	"github.com/hpcloud/tail"
	"github.com/olekukonko/tablewriter"
)

func parseFlags() *ddlog.Config {
	o := ddlog.NewConfig()
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

	// process messages
	metricStore := ddlog.NewMetricStore(config)

	go metricStore.Start(msgChan, func(e *ddlog.MetricEvent) {

		if e.Alert != nil {
			if e.Alert.Done {
				fmt.Printf("[Alert Done] at %s duration=%s\n", e.Alert.End, e.Alert.End.Sub(e.Alert.Start))
			} else {
				fmt.Printf("High traffic generated an alert - hits = %v, triggered at %s", e.Alert.Count, e.Alert.Start)
			}
		}

		fmt.Println("Top 5 Pages")
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Page", "Total Visits"})
		for _, top := range e.TopPages {
			table.Append([]string{top.Name, strconv.Itoa(top.Count)})
		}
		table.Render()
		fmt.Println("")
	})

	done := make(chan bool)
	<-done
}

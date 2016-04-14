package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	"github.com/gee-go/ddlog/ddlog"
	"github.com/gee-go/ddlog/ddlog/cli"
	"github.com/hpcloud/tail"
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

// func initUI(quitChan chan bool) error {
// 	defer ui.Close()
// 	if err := ui.Init(); err != nil {
// 		return err
// 	}

// 	//termui.UseTheme("helloworld")

// 	topKList := ui.NewList()
// 	topKList.Items = []string{"Calculating Top Pages"}
// 	topKList.Height = ui.TermHeight() / 2

// 	ui.Body.AddRows(
// 		ui.NewRow(
// 			ui.NewCol(12, 0, topKList),
// 		),
// 	)

// 	ui.Body.Align()
// 	ui.Render(ui.Body)

// 	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
// 		ui.Render(ui.Body)
// 		ui.Body.Width = ui.TermWidth()
// 		ui.Body.Align()
// 		ui.Render(ui.Body)
// 	})

// 	ui.Handle("/sys/kbd/q", func(ui.Event) {
// 		ui.StopLoop()
// 		ui.Close()
// 		quitChan <- true
// 		close(quitChan)
// 	})
//   ui.Loop()
// 	return nil
// }

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

	ui := cli.NewUI(mon)

	ui.Start()
	// ticker := time.Tick(time.Second)
	// for range ticker {
	// 	fmt.Println(spark.Line(mon.Spark()))
	// }

	// for a := range mon.AlertChan() {
	// 	fmt.Println(a)
	// }
}

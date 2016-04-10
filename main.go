package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/k0kubun/pp"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fn := "/usr/local/var/log/nginx/access.log"

	s := lparse.NewFileScanner()
	go func() {
		<-c
		s.Cleanup()
		os.Exit(1)
	}()

	go s.Tail(fn)

	for l := range s.Line() {
		pp.Println(l)
	}

}

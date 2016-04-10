package main

import (
	"github.com/gee-go/dd_test/src/lparse"
	"github.com/k0kubun/pp"
)

func main() {
	fn := "/usr/local/var/log/nginx/access.log"

	s := lparse.NewFileScanner()
	go s.Tail(fn)

	for l := range s.Line() {
		pp.Println(l)
	}

}

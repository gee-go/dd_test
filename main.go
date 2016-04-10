package main

import (
	"fmt"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/hpcloud/tail"
	"github.com/k0kubun/pp"
)

func main() {
	t, err := tail.TailFile("/usr/local/var/log/nginx/access.log", tail.Config{Follow: true})
	if err != nil {
		fmt.Println(err)
		return
	}

	p := lparse.New(lparse.NewConfig())
	for line := range t.Lines {
		pp.Println(p.Parse(line.Text))
	}
}

package lscan

import "github.com/gee-go/dd_test/src/lparse"

type Scanner interface {
	Line() <-chan *Line
}

type Line struct {
	Msg *lparse.Message
	Err error
}

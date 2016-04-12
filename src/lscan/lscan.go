package lscan

import "github.com/gee-go/dd_test/ddlog"

type Scanner interface {
	Line() <-chan *Line
}

type Line struct {
	Msg *ddlog.Message
	Err error
}

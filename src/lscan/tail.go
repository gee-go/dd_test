package lscan

import (
	"fmt"
	"os"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/hpcloud/tail"
)

type TailScanner struct {
	MsgChan chan *lparse.Message

	t *tail.Tail
	p *lparse.Parser
}

func (s *TailScanner) Cleanup() {
	s.t.Stop()
	s.t.Cleanup()
}

func (s *TailScanner) Start() {
	defer close(s.MsgChan)
	for line := range s.t.Lines {
		if line.Err != nil {
			fmt.Println(line.Err)
			continue
		}

		m, err := s.p.Parse(line.Text)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.MsgChan <- m
	}
}

func Tail(fn string, config *lparse.Config) (*TailScanner, error) {
	t, err := tail.TailFile(fn, tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})

	if err != nil {
		return nil, err
	}

	return &TailScanner{
		MsgChan: make(chan *lparse.Message),
		t:       t,
		p:       lparse.New(config),
	}, nil
}

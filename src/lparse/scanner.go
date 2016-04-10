package lparse

import (
	"os"

	"github.com/hpcloud/tail"
)

type Scanner interface {
	Line() <-chan *Message
}

type FileScanner struct {
	done  chan bool
	lines chan *Line
	err   error
}

func NewFileScanner() *FileScanner {
	return &FileScanner{
		done:  make(chan bool),
		lines: make(chan *Line),
	}
}

func (s *FileScanner) stop() {
	s.done <- true
}

func (s *FileScanner) Line() <-chan *Line {
	return s.lines
}

func (s *FileScanner) Err() error {
	return s.err
}

func (s *FileScanner) Tail(fn string) {
	defer s.stop()

	t, err := tail.TailFile(fn, tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})

	if err != nil {
		s.err = err
		return
	}

	p := New(NewConfig())
	for line := range t.Lines {
		if line.Err != nil {
			s.lines <- &Line{Err: line.Err}
			continue
		}
		m, err := p.Parse(line.Text)
		s.lines <- &Line{Msg: m, Err: err}
	}

	s.err = t.Wait()
	return
}

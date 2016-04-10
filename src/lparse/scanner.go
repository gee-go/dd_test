package lparse

import (
	"os"

	"github.com/hpcloud/tail"
)

type Scanner interface {
	Line() <-chan *Message
}

type FileScanner struct {
	lines chan *Line
	err   error
	tail  *tail.Tail
}

func NewFileScanner() *FileScanner {
	return &FileScanner{
		lines: make(chan *Line),
	}
}

func (s *FileScanner) Cleanup() {
	s.tail.Stop()
	s.tail.Cleanup()
}

func (s *FileScanner) Line() <-chan *Line {
	return s.lines
}

func (s *FileScanner) Err() error {
	return s.err
}

func (s *FileScanner) Tail(fn string) {
	s.tail, s.err = tail.TailFile(fn, tail.Config{
		Follow: true,
		Logger: tail.DiscardingLogger,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: os.SEEK_END,
		},
	})

	if s.err != nil {
		return
	}

	p := New(NewConfig())
	for line := range s.tail.Lines {
		if line.Err != nil {
			s.lines <- &Line{Err: line.Err}
			continue
		}
		m, err := p.Parse(line.Text)
		s.lines <- &Line{Msg: m, Err: err}
	}

	s.err = s.tail.Wait()
	return
}

package lscan

import (
	"os"

	"github.com/gee-go/dd_test/src/lparse"
	"github.com/hpcloud/tail"
)

type TailScanner struct {
}

type FileScanner struct {
	lines  chan *Line
	err    error
	tail   *tail.Tail
	config *lparse.Config
}

func NewFileScanner(config *lparse.Config) *FileScanner {
	return &FileScanner{
		lines:  make(chan *Line),
		config: config,
	}
}

func (s *FileScanner) Cleanup() {
	if s.tail != nil {
		s.tail.Stop()
		s.tail.Cleanup()
	}

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

	p := lparse.New(s.config)
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

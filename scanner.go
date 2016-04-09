package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Scanner struct {
	scanner *bufio.Scanner
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
	}
}

func NewScannerString(s string) *Scanner {
	return NewScanner(strings.NewReader(s))
}

func (p *Scanner) Scan() {
	for p.scanner.Scan() {
		fmt.Println(p.scanner.Text()) // Println will add back the final '\n'
	}
	if err := p.scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

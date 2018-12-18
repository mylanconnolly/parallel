package main

import (
	"bufio"
	"bytes"
	"io"
)

// This is our concurrency-safe line queue, meant to wrap some io.Reader.
type queue struct {
	ch <-chan string
}

func (q queue) readLine() (string, bool) {
	str, more := <-q.ch
	return str, more
}

func newQueue(reader io.Reader, delimMode int) queue {
	ch := make(chan string)

	go func() {
		scanner := newScanner(reader, delimMode)

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return queue{ch: ch}
}

func newScanner(reader io.Reader, mode int) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)

	switch mode {
	case delimNewline:
		scanner.Split(bufio.ScanLines)
	case delimNull:
		scanner.Split(scanNull)
	}
	return scanner
}

// This function is a modification of `bufio.ScanLines` in order to use null
// bytes as a terminator, instead.
func scanNull(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\x00'); i >= 0 {
		// We have a full null-terminated line.
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

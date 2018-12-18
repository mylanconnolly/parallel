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
		scanner := bufio.NewScanner(reader)

		switch delimMode {
		case delimNewline:
			scanner.Split(bufio.ScanLines)
		case delimNull:
			scanner.Split(scanNull)
		}

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return queue{ch: ch}
}

func scanNull(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\x00'); i >= 0 {
		// We have a full null-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

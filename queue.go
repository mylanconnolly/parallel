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

func newQueue(reader io.Reader, splitChar byte, queueBuffer int) queue {
	ch := make(chan string, queueBuffer*2) // Buffer the channel to a reasonable value

	// Build the scanner and start scanning lines into the job queue in the
	// background while we return our new queue.
	go func() {
		scanner := newScanner(reader, splitChar)

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return queue{ch: ch}
}

func newScanner(reader io.Reader, splitChar byte) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)
	scanner.Split(newSplitFunc(splitChar))
	return scanner
}

// This function is used to return a new `bufio.SplitFunc` splitting on
// whichever character the user specifies. The code for this is mostly just
// lifted out of `bufio.ScanLines`, replacing the newline character with a
// paramter.
func newSplitFunc(char byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, char); i >= 0 {
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
}

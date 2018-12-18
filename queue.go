package main

import (
	"bufio"
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

func newQueue(reader io.Reader) queue {
	ch := make(chan string)

	go func() {
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()
	return queue{ch: ch}
}

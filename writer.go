package main

import (
	"bufio"
	"io"
	"sync"
)

// This is used to wrap an io.Writer. It enables us to put a mutex around it,
// which should help ensure that output streams are not getting interspersed.
type writer struct {
	writer io.Writer
	mu     sync.Mutex
}

func newWriter(w io.Writer) *writer {
	return &writer{writer: bufio.NewWriter(w)}
}

// Write is used to implement `io.Writer`
func (w *writer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	n, err = w.writer.Write(p)
	w.mu.Unlock()
	return n, err
}

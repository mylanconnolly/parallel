package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestQueueReadLine(t *testing.T) {
	tests := []struct {
		input     string
		splitChar byte
		want      []string
	}{
		{"one\ttwo", '\t', []string{"one", "two"}},
		{"one\ntwo", '\n', []string{"one", "two"}},
		{"one\rtwo", '\r', []string{"one", "two"}},
		{"abc", 'b', []string{"a", "c"}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			lines := []string{}

			queue := newQueue(strings.NewReader(tt.input), tt.splitChar, 1)

		loop:
			for {
				select {
				case got, ok := <-queue.ch:
					if !ok {
						break loop
					}
					lines = append(lines, got)
				}
			}
			if !reflect.DeepEqual(lines, tt.want) {
				t.Fatalf("Got: %#v Want: %#v", lines, tt.want)
			}
		})
	}
}

package main

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// This creates a stub of a work queue by returning a newline-delimited string
// reader of all of the files in /etc
func testWorkQueue(delim rune) (io.Reader, error) {
	files, err := filepath.Glob("/etc/**")

	if err != nil {
		return nil, err
	}
	return strings.NewReader(strings.Join(files, string('\n'))), nil
}

// This is my generalized benchmark function. It accepts a `*testing.B`, the
// concurrency I want to attempt to work at, the command to run for each item
// in the queue, as well as the args that would be supplied.
//
// This version uses a regular command rather than a template.
func benchmarkCmd(b *testing.B, concurrency int, cmd string, args ...string) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		reader, err := testWorkQueue('\n')

		if err != nil {
			b.Fatal(err)
		}
		pool, err := NewWorkerPool(
			context.Background(),
			ioutil.Discard,
			ioutil.Discard,
			reader,
			byte('\n'),
			cmd,
			"",
			[]string{},
			concurrency,
		)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

		pool.run()
	}
}

// This is my generalized benchmark function. It accepts a `*testing.B`, the
// concurrency I want to attempt to work at, the template to run for each item
// in the queue, as well as the args that would be supplied.
//
// This version uses a template rather than a command
func benchmarkTmpl(b *testing.B, concurrency int, tmpl string, args ...string) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		reader, err := testWorkQueue('\n')

		if err != nil {
			b.Fatal(err)
		}
		pool, err := NewWorkerPool(
			context.Background(),
			ioutil.Discard,
			ioutil.Discard,
			reader,
			byte('\n'),
			"",
			tmpl,
			[]string{},
			concurrency,
		)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()

		pool.run()
	}
}

func BenchmarkXzEtcCmd(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkCmd(b, n, "xz")
		})
	}
}
func BenchmarkCatEtcCmd(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkCmd(b, n, "cat")
		})
	}
}
func BenchmarkEchoEtcCmd(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkCmd(b, n, "echo")
		})
	}
}

func BenchmarkXzEtcTmpl(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkTmpl(b, n, "xz {{.Input}}")
		})
	}
}

func BenchmarkCatEtcTmpl(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkTmpl(b, n, "cat {{.Input}}")
		})
	}
}

func BenchmarkEchoEtcTmpl(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkTmpl(b, n, "echo {{.Input}}")
		})
	}
}

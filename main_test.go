package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var srand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randString(wordLen int) string {
	b := make([]byte, wordLen)
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	for i := range b {
		b[i] = charset[srand.Intn(len(charset))]
	}
	return string(b)
}

// This creates a stub of a work queue by returning a newline-delimited string
// reader of all of the files in /etc
func testWorkQueue(wordLen, queueLen int, delim rune) io.Reader {
	buf := &bytes.Buffer{}

	for i := 0; i < queueLen; i++ {
		if i > 0 {
			buf.WriteRune(delim)
		}
		buf.WriteString(randString(wordLen))
	}
	return buf
}

// This is my generalized benchmark function. It accepts a `*testing.B`, the
// concurrency I want to attempt to work at, the command to run for each item
// in the queue, as well as the args that would be supplied.
//
// This version uses a regular command rather than a template.
func benchmarkCmd(b *testing.B, concurrency, wordLen, queueLen int, cmd string, args ...string) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()

		pool, err := NewWorkerPool(
			context.Background(),
			ioutil.Discard,
			ioutil.Discard,
			testWorkQueue(wordLen, queueLen, '\n'),
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
func benchmarkTmpl(b *testing.B, concurrency, wordLen, queueLen int, tmpl string, args ...string) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()

		pool, err := NewWorkerPool(
			context.Background(),
			ioutil.Discard,
			ioutil.Discard,
			testWorkQueue(wordLen, queueLen, '\n'),
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

func BenchmarkEchoEtcCmd(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkCmd(b, n, 64, 1000, "echo")
		})
	}
}

func BenchmarkEtcTmpl(b *testing.B) {
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			benchmarkTmpl(b, n, 64, 1000, "echo")
		})
	}
}

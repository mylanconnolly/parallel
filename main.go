package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

const (
	delimNewline = 1 << iota
	delimNull
)

func main() {
	jobs := flag.Int("j", runtime.NumCPU(), "Number of jobs to run; defaults to logical CPU core count.")
	argFile := flag.String("a", "", "Use input-file as input source. Only one input source can be specified. In this case, stdin is discarded.")
	nullDelim := flag.Bool("0", false, "Use NUL as delimiter. Normally input lines will end in \\n (newline). If they end in \\0 (NUL), then use this option. It is useful for processing arguments that may contain \\n (newline)")

	flag.Parse()

	if len(os.Args) == 1 || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Must specify program name to run")
		os.Exit(1)
	}
	delimMode := delimNewline

	if *nullDelim {
		delimMode = delimNull
	}
	reader, err := getInput(*argFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Problem getting input source:", err)
		os.Exit(1)
	}
	q := newQueue(reader, delimMode)
	w := worker{queue: q, concurrency: *jobs, program: flag.Arg(0)}

	w.runJobs()
}

func getInput(argFile string) (io.Reader, error) {
	if argFile == "" {
		return os.Stdin, nil
	}
	stat, err := os.Stat(argFile)

	if err != nil || stat.IsDir() {
		return nil, errors.New("must specify a valid file path for -a option")
	}
	return os.Open(argFile)
}

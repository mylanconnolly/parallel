package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
)

func main() {
	jobs := flag.Int("j", runtime.NumCPU(), "Number of jobs to run; defaults to logical CPU core count.")
	argFile := flag.String("a", "", "Use input-file as input source. Only one input source can be specified. In this case, stdin is discarded.")

	flag.Parse()

	if len(os.Args) == 1 || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Must specify program name to run")
		os.Exit(1)
	}
	reader, err := getInput(*argFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Problem getting input source:", err)
		os.Exit(1)
	}
	q := newQueue(reader)
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

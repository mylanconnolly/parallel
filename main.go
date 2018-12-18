package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func main() {
	jobs := flag.Int("j", runtime.NumCPU(), "Number of jobs to run; defaults to logical CPU core count.")
	flag.Parse()

	if len(os.Args) == 1 || len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Must specify program name to run")
		os.Exit(1)
	}
	q := newQueue(os.Stdin)
	w := worker{queue: q, concurrency: *jobs, program: flag.Arg(0)}

	w.runJobs()
}

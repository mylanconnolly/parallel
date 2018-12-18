package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

var jobs = flag.Int("j", runtime.NumCPU(), "The maximum number of jobs to run. By default, the number of logical cores on the local machine.")

func main() {
	flag.Parse()

	program, programArgs, ok := parseArgs(flag.Args())

	if !ok {
		fmt.Fprintln(os.Stderr, "Must specify a command to execute")
		os.Exit(1)
	}
	runJobs(program, programArgs)
}

func parseArgs(args []string) (program string, programArgs []string, ok bool) {
	switch len(args) {
	case 0:
		return "", nil, false
	case 1:
		return args[0], nil, true
	default:
		return args[0], args[1:], true
	}
}

func runJobs(program string, args []string) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, *jobs)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		wg.Add(1)
		sem <- struct{}{}
		input := scanner.Text()

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			if err := runJob(program, input, args); err != nil {
				fmt.Fprintln(os.Stderr, "Error encountered running command:", err.Error())
			}
		}()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	wg.Wait()
}

func runJob(program, input string, args []string) error {
	output, err := exec.Command(program, append(args, input)...).CombinedOutput()
	fmt.Print(string(output))
	return err
}

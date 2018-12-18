package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

// Command-line arguments
var (
	jobs     = flag.Int("j", runtime.NumCPU(), "The maximum number of jobs to run. By default, the number of logical cores on the local machine.")
	argsFile = flag.String("a", "", "Path to args file. If exists, will read lines from file instead of STDIN.")
)

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-j num] [-a path] program", os.Args[0])
		os.Exit(1)
	}
	program, programArgs, ok := parseArgs(flag.Args())

	if !ok {
		fmt.Fprintln(os.Stderr, "Must specify a command to execute")
		os.Exit(1)
	}
	runJobs(program, programArgs, *jobs)
}

// This function pulls double-duty by accepting the positional arguments from
// the command line and splitting it into the program name and any arguments
// for the program. If the input slice is empty, the `ok` value is false, to
// indicate that the user needs to pass a command.
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

// If the input is not empty, assume that it is a file and attempt to open it.
// Otherwise, return os.Stdin as the reader.
func programInput(filePath string) (io.Reader, error) {
	if filePath != "" {
		return os.Open(filePath)
	}
	return os.Stdin, nil
}

// Asynchronously run the program with the given arguments for each of the
// lines of input grabbed from `programInput`.
func runJobs(program string, args []string, concurrency int) {
	reader, err := programInput(*argsFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read from input:", err.Error())
		os.Exit(1)
	}
	runJobsFromReader(program, args, concurrency, reader)
}

func lineReader(reader io.Reader) <-chan string {
	scanner := bufio.NewScanner(reader)
	ch := make(chan string)

	// Set the split function (scan per line)
	scanner.Split(bufio.ScanLines)

	go func() {
		defer close(ch)

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()

	return ch
}

// This is the actual worker used in `runJobs`.
func runJobsFromReader(program string, args []string, concurrency int, reader io.Reader) {
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, concurrency)
	lines := lineReader(reader)

	for line := range lines {
		wg.Add(1)         // Increment waitgroup
		sem <- struct{}{} // Acquire semaphore
		go runJob(program, line, args, sem, &wg)
	}
	wg.Wait()
}

// This function actually executes the program and prints the output from the
// command.
func runJob(program, line string, args []string, sem <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()          // Decrement waitgroup
	defer func() { <-sem }() // Release semaphore

	output, err := exec.Command(program, append(args, line)...).CombinedOutput()

	fmt.Fprint(os.Stdout, string(output))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encountered:", err.Error())
	}
}

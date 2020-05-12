package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
)

const (
	delimNewLine = byte('\n')
	delimNull    = byte('\x00')
)

func main() {
	nullDelim := flag.Bool("0", false, "Use NUL as delimiter. Normally input lines will end in \\n (newline). If they end in \\0 (NUL), then use this option. It is useful for processing arguments that may contain \\n (newline)")
	argFile := flag.String("a", "", "Use input-file as input source. Only one input source can be specified. In this case, stdin is discarded.")
	jobs := flag.Int("j", runtime.NumCPU(), "Number of jobs to run; defaults to logical CPU core count.")
	template := flag.String("t", "", "Specify a command template, which is used to override the default behavior of one command per line, with the line appended to the command and any positional arguments that exist")

	flag.Parse()

	args := flag.Args()

	if *template == "" && len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Must specify command to run or a template using the -t flag\n")
		os.Exit(1)
	}
	delimMode := delimNewLine

	if *nullDelim {
		delimMode = delimNull
	}
	reader, err := getInput(*argFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem getting input source: %s\n", err)
		os.Exit(1)
	}
	var programArgs []string

	if len(args) > 1 {
		programArgs = args[1:]
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	cmd := ""

	if len(args) > 0 {
		cmd = args[0]
	}
	w, err := NewWorkerPool(
		ctx,
		os.Stdout,
		os.Stderr,
		reader,
		delimMode,
		cmd,
		*template,
		programArgs,
		*jobs,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		cancel()
		fmt.Fprintf(os.Stderr, "Caught SIGINT, exiting gracefully (send once more to exit immediately)\n")
		<-c
		os.Exit(130)
	}()
	w.run()
}

func getInput(argFile string) (io.Reader, error) {
	if argFile == "" {
		return os.Stdin, nil
	}
	stat, err := os.Stat(argFile)

	if err != nil || stat.IsDir() {
		return nil, fmt.Errorf("must specify a valid file path for -a option")
	}
	return os.Open(argFile)
}

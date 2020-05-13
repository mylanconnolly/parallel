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
	nullDelim := flag.Bool(
		"0",
		false,
		"Use NUL as delimiter. Normally input lines will end in \\n (newline). If they end in \\0 (NUL), then use this option. It is useful for processing arguments that may contain \\n (newline)",
	)
	argFile := flag.String(
		"a",
		"",
		"Use input-file as input source. Only one input source can be specified. In this case, stdin is discarded.",
	)
	jobs := flag.Int(
		"j",
		runtime.NumCPU(),
		"Number of jobs to run; defaults to logical CPU core count.",
	)
	template := flag.String(
		"t",
		"",
		"Specify a command template, which is used to override the default behavior of one command per line, with the line appended to the command and any positional arguments that exist",
	)
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
	cmd := getCmd(args)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if cmd == "" && *template == "" {
		fmt.Fprintln(os.Stderr, "Must provide a command or command template")
		os.Exit(1)
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
	// Check for SIGINT and attempt to tell all jobs to cancel. If the user sends
	// another SIGINT then just forcefully quit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		cancel()
		fmt.Fprintf(
			os.Stderr,
			"Caught SIGINT, exiting gracefully (send once more to exit immediately)\n",
		)
		<-c
		os.Exit(130)
	}()

	// Start the worker process.
	w.run()
}

// If the args slice has at least one value then the first value would be the
// command to run. If there isn't anything in it, then just return an empty
// string. Presumably the template will be provided if there is no command.
func getCmd(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
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

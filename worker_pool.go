package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type WorkerPool struct {
	args        []string
	cmd         string
	concurrency int
	ctx         context.Context
	err         *writer
	out         *writer
	queue       queue
	runner      func(line string)
	start       time.Time
	template    *template.Template
}

// NewWorkerPool is designed to set up the worker pool for use.
func NewWorkerPool(
	ctx context.Context,
	stdout, stderr io.Writer,
	reader io.Reader,
	splitChar byte,
	cmd, tmpl string,
	args []string,
	concurrency int,
) (*WorkerPool, error) {
	var (
		parsedTmpl *template.Template
		path       string
		err        error
	)
	if cmd != "" {
		path, err = exec.LookPath(cmd)

		if err != nil {
			return nil, err
		}
	}
	w := &WorkerPool{
		args:        args,
		cmd:         path,
		concurrency: concurrency,
		ctx:         ctx,
		err:         newWriter(stderr),
		out:         newWriter(stdout),
		queue:       newQueue(reader, splitChar, concurrency),
		start:       time.Now(),
		template:    parsedTmpl,
	}
	if tmpl != "" {
		parsedTmpl, err = template.New("cmd").Funcs(tmplFuncs).Parse(tmpl)

		if err != nil {
			return nil, err
		}
		w.template = parsedTmpl
		w.runner = w.runCmdTemplate
	} else {
		w.runner = w.runCmd
	}
	return w, nil
}

func (w *WorkerPool) run() {
	wg := sync.WaitGroup{}

	for i := 0; i < w.concurrency; i++ {
		wg.Add(1)
		go w.startWorker(&wg)
	}
	wg.Wait()
	w.out.Flush()
	w.err.Flush()
}

func (w *WorkerPool) startWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Check if the user cancelled the run
		select {
		case <-w.ctx.Done():
			return
		case line, open := <-w.queue.ch:
			if !open {
				return
			}
			w.runner(line)
		default:
		}
	}
}

func (w *WorkerPool) runCmd(input string) {
	args := append(w.args, input)
	cmd := newCmd(w.out, w.err, w.cmd, args...)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(w.err, "Failed to run command: `%s %s` %s\n", w.cmd, strings.Join(args, " "), err)
	}
}

func (w *WorkerPool) runCmdTemplate(input string) {
	buf := bytes.Buffer{}
	ctx := Ctx{
		Cmd:   w.cmd,
		Input: input,
		Start: w.start,
		Time:  time.Now(),
	}
	if err := w.template.Execute(&buf, ctx); err != nil {
		fmt.Fprintf(w.err, "%s\n", err)
		return
	}
	words := shellParser(buf.String())
	cmd := newCmd(w.out, w.err, words[0], words[1:]...)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(w.err, "Failed to run command `%s`: %s\n", buf.String(), err)
	}
}

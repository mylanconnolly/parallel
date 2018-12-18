package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type worker struct {
	concurrency int
	queue       queue
	program     string
	args        []string
}

func (w worker) runJobs() {
	wg := sync.WaitGroup{}

	for i := 0; i < w.concurrency; i++ {
		wg.Add(1)
		go w.runJob(&wg)
	}
	wg.Wait()
}

func (w worker) runJob(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		line, ok := w.queue.readLine()

		if !ok {
			break
		}
		out, err := exec.Command(w.program, append(w.args, line)...).CombinedOutput()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Stdout.Write(out)
	}
}

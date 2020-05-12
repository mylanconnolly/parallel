// +build windows

package main

import (
	"io"
	"os/exec"
)

func newCmd(stdout, stderr io.Writer, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	return c
}

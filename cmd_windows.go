// +build windows

package main

import (
	"io"
	"os/exec"
)

// Build the command. This is the Windows version of the function which does not
// allow us to prevent propagating SIGINT to child processes. For this reason,
// we cannot gracefully quit when compiled for Windows.
func newCmd(stdout, stderr io.Writer, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	return c
}

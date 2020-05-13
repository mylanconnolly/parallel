// +build darwin freebsd linux

package main

import (
	"io"
	"os/exec"
	"syscall"
)

// Build the command. This is the UNIX version of the function, which makes an
// effort to prevent SIGINT from propagating directly to child processes. This
// gives the app the ability to gracefully quit as the workers finish their
// current jobs.
func newCmd(stdout, stderr io.Writer, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return c
}

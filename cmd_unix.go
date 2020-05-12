// +build darwin freebsd linux

package main

import (
	"io"
	"os/exec"
	"syscall"
)

func newCmd(stdout, stderr io.Writer, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return c
}

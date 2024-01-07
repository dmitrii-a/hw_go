package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

const (
	OK = iota
	Fail
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for name, e := range env {
		var err error
		if e.NeedRemove {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, e.Value)
		}
		if err != nil {
			return Fail
		}
	}
	//nolint:gosec
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		log.Printf("Run command error: %s", err)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return Fail
	}
	return OK
}

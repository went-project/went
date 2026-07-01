//go:build windows

package handlers

import (
	"os"
	"os/exec"
)

func setupProcessGroup(cmd *exec.Cmd) {}

func signalProcessGroup(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func killProcessGroup(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func terminationSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}

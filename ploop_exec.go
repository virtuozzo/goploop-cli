package ploop

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func ploopRunCmd(stdout io.Writer, args ...string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("ploop", args...)
	cmd.Stdout = stdout
	cmd.Stderr = &stderr

	fmt.Printf("Run: %s\n", strings.Join([]string{cmd.Path, strings.Join(cmd.Args[1:], " ")}, " "))

	err := cmd.Run()
	if err == nil {
		return nil
	}

	// Command returned an error, get the stderr
	errStr := stderr.String()
	// Get the exit code (Unix-specific)
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			errCode := status.ExitStatus()
			return &Err{c: errCode, s: errStr}
		}
	}
	// unknown exit code
	return &Err{c: -1, s: errStr}
}

func ploop(args ...string) error {
	return ploopRunCmd(os.Stdout, args...)
}

func ploopOut(args ...string) (string, error) {
	var stdout bytes.Buffer
	ret := ploopRunCmd(&stdout, args...)
	return stdout.String(), ret
}

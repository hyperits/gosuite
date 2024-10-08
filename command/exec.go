package command

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecuteCommand runs a command and returns the output or an error.
func ExecuteCommand(command string, args ...string) (string, error) {
	// Create the command
	cmd := exec.Command(command, args...)

	// Capture the output and error
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %v - %s", err, stderr.String())
	}

	// Return the output
	return out.String(), nil
}

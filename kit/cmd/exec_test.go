package cmd_test

import (
	"testing"

	"github.com/hyperits/gosuite/kit/cmd"
)

func TestExecuteCommand(t *testing.T) {

	output, err := cmd.ExecuteCommand("ls", "-l")
	if err != nil {
		t.Error("Error executing command", err)
	} else {
		t.Log("Command output:", output)
	}

}

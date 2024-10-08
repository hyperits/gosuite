package command_test

import (
	"testing"

	"github.com/hyperits/gosuite/command"
)

func TestExecuteCommand(t *testing.T) {

	output, err := command.ExecuteCommand("ls", "-l")
	if err != nil {
		t.Error("Error executing command", err)
	} else {
		t.Log("Command output:", output)
	}

}

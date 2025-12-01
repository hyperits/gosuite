package cmd_test

import (
	"context"
	"testing"
	"time"

	"github.com/hyperits/gosuite/kit/cmd"
)

func TestExecuteCommand(t *testing.T) {
	output, err := cmd.ExecuteCommand("echo", "hello")
	if err != nil {
		t.Error("Error executing command", err)
	} else {
		t.Log("Command output:", output)
	}
}

func TestExecuteCommandWithContext(t *testing.T) {
	ctx := context.Background()
	output, err := cmd.ExecuteCommandWithContext(ctx, "echo", "hello with context")
	if err != nil {
		t.Error("Error executing command with context", err)
	} else {
		t.Log("Command output:", output)
	}
}

func TestExecuteCommandWithTimeout(t *testing.T) {
	output, err := cmd.ExecuteCommandWithTimeout(5*time.Second, "echo", "hello with timeout")
	if err != nil {
		t.Error("Error executing command with timeout", err)
	} else {
		t.Log("Command output:", output)
	}
}

func TestExecuteCommandResult(t *testing.T) {
	ctx := context.Background()
	result, err := cmd.ExecuteCommandResult(ctx, "echo", "hello result")
	if err != nil {
		t.Error("Error executing command result", err)
	} else {
		t.Logf("Command result: stdout=%s, stderr=%s, exitCode=%d", result.Stdout, result.Stderr, result.ExitCode)
	}
}

func TestExecuteCommandTimeout(t *testing.T) {
	// 测试超时场景
	_, err := cmd.ExecuteCommandWithTimeout(1*time.Millisecond, "sleep", "10")
	if err == nil {
		t.Error("Expected timeout error, got nil")
	} else {
		t.Log("Expected timeout error:", err)
	}
}

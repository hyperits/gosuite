package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Result 命令执行结果
type Result struct {
	Stdout   string // 标准输出
	Stderr   string // 错误输出
	ExitCode int    // 退出码
}

// ExecuteCommand 执行命令并返回输出或错误
func ExecuteCommand(command string, args ...string) (string, error) {
	return ExecuteCommandWithContext(context.Background(), command, args...)
}

// ExecuteCommandWithContext 使用指定上下文执行命令
func ExecuteCommandWithContext(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// 检查是否是上下文取消或超时
		if ctx.Err() != nil {
			return "", fmt.Errorf("command cancelled or timed out: %w", ctx.Err())
		}
		return "", fmt.Errorf("command failed: %v - %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// ExecuteCommandWithTimeout 使用超时执行命令
func ExecuteCommandWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ExecuteCommandWithContext(ctx, command, args...)
}

// ExecuteCommandResult 执行命令并返回详细结果
func ExecuteCommandResult(ctx context.Context, command string, args ...string) (*Result, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		// 检查是否是上下文取消或超时
		if ctx.Err() != nil {
			return result, fmt.Errorf("command cancelled or timed out: %w", ctx.Err())
		}

		// 尝试获取退出码
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result, fmt.Errorf("command failed: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

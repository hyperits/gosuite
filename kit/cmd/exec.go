package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
)

// ExecuteCommand 执行命令并返回输出或错误
func ExecuteCommand(command string, args ...string) (string, error) {
	// 创建命令
	cmd := exec.Command(command, args...)

	// 捕获标准输出和错误输出
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command failed: %v - %s", err, stderr.String())
	}

	// 返回输出结果
	return out.String(), nil
}

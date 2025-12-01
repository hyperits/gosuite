package debug

import (
	"errors"
	"runtime"
	"strings"
)

// RuntimeInfo 运行时信息结构体
type RuntimeInfo struct {
	File     string // 文件路径
	Line     int    // 行号
	Function string // 函数名（不含包名）
	Err      error  // 错误信息
}

// GetCurrentFunctionInfo 获取当前函数的运行时信息
func GetCurrentFunctionInfo() *RuntimeInfo {
	res := &RuntimeInfo{
		File:     "init",
		Line:     0,
		Function: "unknown",
		Err:      nil,
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		res.Err = errors.New("failed to get function info")
		return res
	}

	function := runtime.FuncForPC(pc).Name()
	parts := strings.Split(function, ".")
	if len(parts) > 1 {
		function = parts[len(parts)-1]
	}

	res.File = file
	res.Line = line
	res.Function = function
	return res
}

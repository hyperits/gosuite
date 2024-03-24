package debugger

import (
	"errors"
	"runtime"
	"strings"
)

type RuntimeInfo struct {
	File     string
	Line     int
	Function string // Function without package
	Err      error
}

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

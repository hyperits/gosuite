package debugger

import (
	"errors"
	"runtime"
)

type RuntimeInfo struct {
	File     string
	Line     int
	Function string
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
	res.File = file
	res.Line = line
	res.Function = function
	return res
}

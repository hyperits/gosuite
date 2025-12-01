package logger_test

import (
	"testing"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/hyperits/gosuite/logger"
)

func TestLogLevels(t *testing.T) {
	// 测试基础日志方法
	logger.Debugf("hello, %s!", "world")
	logger.Infof("hello, %s!", "world")
	logger.Warnf("hello, %s!", "world")
	logger.Errorf("hello, %s!", "world")
}

func TestLogWithRuntimeInfo(t *testing.T) {
	// 测试带运行时信息的日志方法
	logger.DebugRTf(debug.GetCurrentFunctionInfo(), "debug with runtime info")
	logger.InfoRTf(debug.GetCurrentFunctionInfo(), "info with runtime info")
	logger.WarnRTf(debug.GetCurrentFunctionInfo(), "warn with runtime info")
	logger.ErrorRTf(debug.GetCurrentFunctionInfo(), "error with runtime info")
}

func TestLogInit(t *testing.T) {
	// 测试配置初始化
	cfg := logger.DefaultConfig()
	cfg.Level = logger.DebugLevel
	logger.Init(cfg)

	logger.Debugf("after init: debug message should appear")
	logger.Infof("after init: info message")
}

func TestLogSetLevel(t *testing.T) {
	// 测试设置日志级别
	if err := logger.SetStrLevel("debug"); err != nil {
		t.Errorf("SetStrLevel failed: %v", err)
	}

	if err := logger.SetStrLevel("invalid"); err == nil {
		t.Error("SetStrLevel should return error for invalid level")
	}
}

// 注意：Fatalf 和 Panicf 会导致程序退出，不适合在单元测试中调用
// func TestLogFatal(t *testing.T) {
// 	logger.Fatalf("this would exit the program")
// }
// func TestLogPanic(t *testing.T) {
// 	logger.Panicf("this would panic")
// }


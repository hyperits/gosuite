package log_test

import (
	"testing"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/hyperits/gosuite/log"
)

func TestLogLevels(t *testing.T) {
	// 测试基础日志方法
	log.Debugf("hello, %s!", "world")
	log.Infof("hello, %s!", "world")
	log.Warnf("hello, %s!", "world")
	log.Errorf("hello, %s!", "world")
}

func TestLogWithRuntimeInfo(t *testing.T) {
	// 测试带运行时信息的日志方法
	log.DebugRTf(debug.GetCurrentFunctionInfo(), "debug with runtime info")
	log.InfoRTf(debug.GetCurrentFunctionInfo(), "info with runtime info")
	log.WarnRTf(debug.GetCurrentFunctionInfo(), "warn with runtime info")
	log.ErrorRTf(debug.GetCurrentFunctionInfo(), "error with runtime info")
}

func TestLogInit(t *testing.T) {
	// 测试配置初始化
	cfg := log.DefaultConfig()
	cfg.Level = log.DebugLevel
	log.Init(cfg)

	log.Debugf("after init: debug message should appear")
	log.Infof("after init: info message")
}

func TestLogSetLevel(t *testing.T) {
	// 测试设置日志级别
	if err := log.SetStrLevel("debug"); err != nil {
		t.Errorf("SetStrLevel failed: %v", err)
	}

	if err := log.SetStrLevel("invalid"); err == nil {
		t.Error("SetStrLevel should return error for invalid level")
	}
}

// 注意：Fatalf 和 Panicf 会导致程序退出，不适合在单元测试中调用
// func TestLogFatal(t *testing.T) {
// 	log.Fatalf("this would exit the program")
// }
// func TestLogPanic(t *testing.T) {
// 	log.Panicf("this would panic")
// }

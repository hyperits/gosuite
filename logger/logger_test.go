package logger_test

import (
	"testing"

	"github.com/hyperits/gosuite/logger"
)

func TestLogDebug(t *testing.T) {
	logger.Debugf("hello, %s!", "world")
	logger.Infof("hello, %s!", "world")
	logger.Warnf("hello, %s!", "world")
	logger.Errorf("hello, %s!", "world")
	logger.Fatalf("hello, %s!", "world")
	logger.Panicf("hello, %s!", "world")
}

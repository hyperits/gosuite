package log_test

import (
	"testing"

	"github.com/hyperits/gosuite/kit/debug"
	"github.com/hyperits/gosuite/log"
)

func TestLogDebug(t *testing.T) {
	log.Debugf("hello, %s!", "world")
	log.Infof("hello, %s!", "world")
	log.Warnf("hello, %s!", "world")
	log.Errorf("hello, %s!", "world")
	log.InfoRTf(debug.GetCurrentFunctionInfo(), "hello, %s!", "world")
	log.Fatalf("hello, %s!", "world")
	log.Panicf("hello, %s!", "world")
}

package sms_test

import (
	"testing"

	"github.com/hyperits/gosuite/component/sms"
	"github.com/hyperits/gosuite/config"
)

func TestAliyunSms(t *testing.T) {
	conf := config.LoadConfig("../../../configs/config.yaml")
	comp := sms.NewSmsComp(&conf.Sms)

	err := comp.SendByAliyun("18888888888", "123456")
	if err != nil {
		t.Fatalf("Error send aliyun sms: %v", err)
	}
}

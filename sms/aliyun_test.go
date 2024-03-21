package sms_test

import (
	"testing"

	"github.com/hyperits/gosuite/sms"
)

func TestAliyunSms(t *testing.T) {
	conf := sms.SmsConfig{}
	comp := sms.NewSmsComponent(&conf)

	err := comp.SendByAliyun("18888888888", "123456")
	if err != nil {
		t.Fatalf("Error send aliyun sms: %v", err)
	}
}

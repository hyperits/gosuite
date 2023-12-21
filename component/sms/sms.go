package sms

import (
	"github.com/hyperits/gosuite/config"
)

type SmsComp struct {
	conf *config.SmsConfig
}

func NewSmsComp(conf *config.SmsConfig) *SmsComp {
	return &SmsComp{
		conf: conf,
	}
}

func (c *SmsComp) GetProvider() string {
	return c.conf.Provider
}

package sms

type SmsComponent struct {
	conf *SmsConfig
}

func NewSmsComponent(conf *SmsConfig) *SmsComponent {
	return &SmsComponent{
		conf: conf,
	}
}

func (c *SmsComponent) GetProvider() string {
	return c.conf.Provider
}

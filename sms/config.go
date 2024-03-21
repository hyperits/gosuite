package sms

type SmsConfig struct {
	Provider string          `yaml:"provider"` // aliyun | etc.
	Aliyun   AliyunSmsConfig `yaml:"aliyun,omitempty"`
}

type AliyunSmsConfig struct {
	Region       string `yaml:"region"`
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	SignName     string `yaml:"sign_name"`
	TemplateCode string `yaml:"template_code"`
}

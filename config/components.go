package config

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
}

type S3Config struct {
	Endpoint       string `yaml:"endpoint"`
	AccessKey      string `yaml:"access_key"`
	Secret         string `yaml:"secret"`
	Bucket         string `yaml:"bucket"`
	Region         string `yaml:"region"`
	Secure         bool   `yaml:"secure"`
	ForcePathStyle bool   `yaml:"force_path_style"`
}

type MailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

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

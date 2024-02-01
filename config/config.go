package config

import (
	"fmt"
	"strings"

	serviceRedis "github.com/hyperits/gosuite/store/redis"
	"gopkg.in/yaml.v3"
)

type SuiteConfig struct {
	Address     string                   `yaml:"address"`               // bind address | :1234
	LogLevel    string                   `yaml:"log_level,omitempty"`   // log | error/warn/info
	MySQL       MySQLConfig              `yaml:"mysql"`                 // store
	S3          S3Config                 `yaml:"s3"`                    // s3
	Redis       serviceRedis.RedisConfig `yaml:"redis,omitempty"`       // redis
	Mail        MailConfig               `yaml:"mail,omitempty"`        // mail
	Sms         SmsConfig                `yaml:"sms,omitempty"`         // sms
	Development bool                     `yaml:"development,omitempty"` // deployment mode
}

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

func NewSuiteConfig(confString string, strictMode bool) (*SuiteConfig, error) {
	// start with defaults
	conf := &SuiteConfig{
		Address:  ":50051",
		LogLevel: "debug",
		MySQL:    MySQLConfig{},
		S3:       S3Config{},
		Redis:    serviceRedis.RedisConfig{},
	}

	if confString != "" {
		decoder := yaml.NewDecoder(strings.NewReader(confString))
		decoder.KnownFields(strictMode)
		if err := decoder.Decode(conf); err != nil {
			return nil, fmt.Errorf("could not parse config: %v", err)
		}
	}

	if conf.LogLevel == "" && conf.Development {
		conf.LogLevel = "debug"
	}

	return conf, nil
}

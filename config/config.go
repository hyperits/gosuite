package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	serviceRedis "github.com/hyperits/gosuite/store/redis"
	"gopkg.in/yaml.v3"
)

// ----------------------------------------------------------------------------
// Example Config, You can impl your own like:
// 	type AppConfig struct {
// 		Address string      `yaml:"address"` // bind address | :1234
//		MySQL   MySQLConfig `yaml:"mysql"`   // store
// 	}
// ----------------------------------------------------------------------------

type SuiteConfig struct {
	MySQL       MySQLConfig              `yaml:"mysql"`                 // store
	S3          S3Config                 `yaml:"s3"`                    // s3
	Redis       serviceRedis.RedisConfig `yaml:"redis,omitempty"`       // redis
	Mail        MailConfig               `yaml:"mail,omitempty"`        // mail
	Sms         SmsConfig                `yaml:"sms,omitempty"`         // sms
	Development bool                     `yaml:"development,omitempty"` // deployment mode
}

func NewSuiteConfig(confString string, strictMode bool) (*SuiteConfig, error) {
	// start with defaults
	conf := &SuiteConfig{
		MySQL: MySQLConfig{},
		S3:    S3Config{},
		Redis: serviceRedis.RedisConfig{},
	}

	if confString != "" {
		decoder := yaml.NewDecoder(strings.NewReader(confString))
		decoder.KnownFields(strictMode)
		if err := decoder.Decode(conf); err != nil {
			return nil, fmt.Errorf("could not parse config: %v", err)
		}
	}

	return conf, nil
}

func AutoLoadSuiteConfig() *SuiteConfig {
	return LoadSuiteConfig("")
}

func LoadSuiteConfig(location string) *SuiteConfig {
	configPath := ""
	// param first
	if location != "" {
		configPath = location
	} else {
		// env second
		envValue := os.Getenv("CONFIG_PATH")
		if envValue == "" {
			// debug default
			configPath = "configs/config.yaml"
		}
	}

	configBody, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config content of %v: %v", configPath, err)
	}
	conf, err := NewSuiteConfig(string(configBody), true)
	if err != nil {
		log.Fatalf("Failed to init config from %v: %v", configPath, err)
	}
	return conf
}

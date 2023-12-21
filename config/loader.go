package config

import (
	"log"
	"os"
)

func AutoLoadConfig() *Config {
	return LoadConfig("")
}

func LoadConfig(location string) *Config {
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
	conf, err := NewConfig(string(configBody), true)
	if err != nil {
		log.Fatalf("Failed to init config from %v: %v", configPath, err)
	}
	return conf
}

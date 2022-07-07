package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Application Application
	Cache       Cache
	Database    Database
	Kafka       Kafka
	Tracing     Tracing
}

func ParseFile(fileBytes []byte) (*File, error) {
	configFile := File{}

	err := yaml.Unmarshal(fileBytes, &configFile)
	if err != nil {
		return nil, err
	}

	return &configFile, nil
}

func ParseConfig(filepath string) (*Config, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	configFile, err := ParseFile(b)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		Application: configFile.Application,
		Cache:       configFile.Cache,
		Database:    configFile.Database,
		Kafka:       configFile.Kafka,
		Tracing:     configFile.Tracing,
	}

	return &cfg, nil
}

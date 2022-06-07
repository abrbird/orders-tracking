package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Application Application
	Database    Database
	Kafka       Kafka
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
		Database:    configFile.Database,
		Kafka:       configFile.Kafka,
	}

	return &cfg, nil
}

package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        string `envconfig:"PORT" default:"8080"`
	ServiceName string `envconfig:"SERVICE_NAME"`
	InstanceId  string `envconfig:"INSTANCE_ID"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	URLList     string `default:"example_list.txt"`
	Concurrency int    `default:"3"`
}

func NewConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("envconfig.Process: %w", err)
	}

	return cfg, nil
}

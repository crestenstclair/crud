package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	DYNAMODB_TABLE  string `env:"DYNAMODB_TABLE,required"`
	RequestTimeoutMS int    `env:"REQUEST_TIMEOUT_MS" envDefault:"200"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

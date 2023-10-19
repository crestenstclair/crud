package config

import (
	"github.com/caarlos0/env/v9"
)

type Config struct {
	DynamoTableName string `env:"DYNAMO_TABLE_NAME,required"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

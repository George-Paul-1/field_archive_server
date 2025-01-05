package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DB_Url string `env:"DATABASE_URL,required"`
	Port   string `env:"PORT,required"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	err = env.Parse(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

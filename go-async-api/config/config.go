package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabaesPort     string `env:"DB_PORT"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
}

func New() (*Config, error) {
	var cfg Config
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) DatabaseUrl() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabaesPort,
		c.DatabaseName,
	)
}

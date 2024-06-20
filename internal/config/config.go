package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RepositoryType  string        `env:"REPO_TYPE" envDefault:"postgres"`
	Port            string        `env:"PORT" envDefault:"8080"`
	DBConn          string        `env:"DBConn" envDefault:"postgres://postgres:postgres@db/postgres?sslmode=disable"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
}

func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

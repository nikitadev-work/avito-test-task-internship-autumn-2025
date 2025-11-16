package config

import (
	"fmt"

	env "github.com/caarlos0/env/v10"
)

// Config holds application configuration loaded from env
type Config struct {
	App        App
	Log        Log
	HTTP       HTTP
	PostgreSQL PostgreSQL
}

type App struct {
	Name    string `env:"APP_NAME,required"`
	Version string `env:"APP_VERSION,required"`
}

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
}

type HTTP struct {
	Port string `env:"HTTP_PORT,required"`
}

type PostgreSQL struct {
	User       string `env:"DB_USER,required"`
	Password   string `env:"DB_PASSWORD,required"`
	Host       string `env:"DB_HOST,required"`
	Port       string `env:"DB_PORT,required"`
	Name       string `env:"DB_NAME,required"`
	SslEnabled bool   `env:"DB_SSL_ENABLED,required"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}

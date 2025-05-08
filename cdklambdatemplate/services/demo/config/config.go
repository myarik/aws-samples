package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	// Config -.
	Config struct {
		App         App
		Log         Log
		Environment Env
	}

	// App -.
	App struct {
		Version string `env:"APP_VERSION,required"`
	}
	// Log -.
	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	Env struct {
		Enabled string `env:"Environment" envDefault:"staging"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

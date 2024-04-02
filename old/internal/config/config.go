package config

import (
	"context"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/logconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/serverconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/storageconfig"
	"github.com/sethvargo/go-envconfig"
)

type AppConfig struct {
	Storage storageconfig.StorageConfig `env:", prefix=STORE_"`
	Server  serverconfig.ServerConfig   `env:", prefix=SERVER_"`
	Log     logconfig.LogConfig         `env:", prefix=LOG_"`

	ModeStr string `env:"MODE, default=DEBUG"`
	Mode    Mode
}

func (c *AppConfig) Close(ctx context.Context) error {
	if err := c.Storage.Close(ctx); err != nil {
		return fmt.Errorf("closing storage config: %w", err)
	}
	return nil
}

type Mode int

const (
	Debug Mode = iota
	Release
)

// Returns a new config and a setup viper instance if there is no error;
// If encounters an error, returns all nil except the error
func New(ctx context.Context) (*AppConfig, error) {
	conf := AppConfig{}

	if err := envconfig.Process(ctx, &conf); err != nil {
		return nil, err
	}

	switch conf.ModeStr {
	case "DEBUG":
		conf.Mode = Debug
	case "RELEASE":
		conf.Mode = Release
	default:
		return nil, fmt.Errorf("application mode %s is not allowed", conf.ModeStr)
	}

	return &conf, nil
}

// Runs configurations for components of the app
func (c *AppConfig) Configure(ctx context.Context, logger applog.Logger) error {
	if err := c.Storage.Configure(ctx, logger); err != nil {
		return fmt.Errorf("configuring storage: %w", err)
	}

	return nil
}

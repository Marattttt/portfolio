package config

import (
	"context"
	"fmt"

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

	if err := conf.Storage.Configure(); err != nil {
		return nil, fmt.Errorf("configuring storage: %w", err)
	}
	// var (
	// 	server serverconfig.ServerConfig
	// 	db     storageconfig.StorageConfig
	// 	log    logconfig.LogConfig
	// )
	//
	// if err := envconfig.Process(ctx, &server); err != nil {
	// 	return nil, fmt.Errorf("creating server config: %w", err)
	// }
	// conf.Server = server
	//
	// if err := envconfig.Process(ctx, &db); err != nil {
	// 	return nil, fmt.Errorf("creating storage config: %w", err)
	// }
	// conf.Storage = db
	//
	// if err := envconfig.Process(ctx, &log); err != nil {
	// 	return nil, fmt.Errorf("creating log config: %w", err)
	// }
	// conf.Log = log
	//

	return &conf, nil
}

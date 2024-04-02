package storageconfig

import (
	"context"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
)

type StorageConfig struct {
	DB *dbconfig.DbConfig `env:", prefix=DB_"`
}

func (c *StorageConfig) Configure(ctx context.Context, logger applog.Logger) error {
	if err := c.DB.Configure(); err != nil {
		return fmt.Errorf("configuring db: %w", err)
	}

	if err := c.DB.Migrate(ctx, logger); err != nil {
		return fmt.Errorf("applying migrations: %w", err)
	}
	return nil

}

func (c *StorageConfig) Close(ctx context.Context) error {
	var err error
	if c.DB != nil {
		err = c.DB.Close(ctx)
	}
	return err
}

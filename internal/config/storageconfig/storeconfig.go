package storageconfig

import (
	"context"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
)

type StorageConfig struct {
	DB *dbconfig.DbConfig `env:", prefix=DB_"`
}

func (c *StorageConfig) Close(ctx context.Context) error {
	var err error
	if c.DB != nil {
		err = c.DB.Close(ctx)
	}
	return err
}

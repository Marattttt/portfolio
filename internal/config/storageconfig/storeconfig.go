package storageconfig

import (
	"context"
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
	"github.com/spf13/viper"
)

type StorageConfig struct {
	DB *dbconfig.DbConfig
}

func (c *StorageConfig) Close(_ context.Context) error {
	var err error
	if c.DB != nil {
		err = c.DB.Close()
	}
	return err
}

func New(vpr *viper.Viper) (*StorageConfig, error) {
	var conf StorageConfig
	dbconf, err := dbconfig.New(vpr)
	if err != nil {
		return nil, fmt.Errorf("creating db config: %w", err)
	}
	conf.DB = dbconf
	return &conf, nil
}

package storageconfig

import (
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/dbconfig"
	"github.com/spf13/viper"
)

type StorageConfig struct {
	DB *dbconfig.DbConfig
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

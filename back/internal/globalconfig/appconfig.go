package appconfig

import (
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/dbconfig"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	DbConf      dbconfig.DbConfig
	ServeerConf ServerConfig
}

type ServerConfig struct {
	ListeningAdress string
}

func CreateGlobalConfig() (*GlobalConfig, viper.Viper, error) {

}

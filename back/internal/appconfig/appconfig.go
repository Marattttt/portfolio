package appconfig

import (
	"errors"
	"strconv"

	"github.com/Marattttt/portfolio/portfolio_back/internal/db/dbconfig"
	"github.com/spf13/viper"
)

type AppConfig struct {
	DB     dbconfig.DbConfig
	Server ServerConfig
}

type ServerConfig struct {
	ListenOn string
}

func CreateAppConfig() (*AppConfig, *viper.Viper, error) {
	appConf := AppConfig{}
	vpr := viper.New()
	vpr.SetEnvPrefix("PORTFOLIO")

	if serverConf, err := SetupServerConfig(*vpr); err != nil {
		return nil, nil, err
	} else {
		appConf.Server = *serverConf
	}

	if dbConf, err := dbconfig.CreateConfig(*vpr); err != nil {
		return nil, nil, err
	} else {
		appConf.DB = *dbConf
	}

	return &appConf, vpr, nil
}

func SetupServerConfig(vpr viper.Viper) (*ServerConfig, error) {
	vpr.BindEnv("PORT")
	listenOn := vpr.GetString("PORT")

	if listenOn == "" {
		return nil, errors.New("Cannot get PORT viper key")
	}

	if _, err := strconv.Atoi(listenOn); err == nil {
		listenOn = ":" + listenOn
	}

	return &ServerConfig{
		ListenOn: listenOn,
	}, nil
}

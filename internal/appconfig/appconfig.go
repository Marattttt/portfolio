package appconfig

import (
	"errors"
	"strconv"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/db/dbconfig"
	"github.com/spf13/viper"
)

type AppConfig struct {
	DB     dbconfig.DbConfig
	Server ServerConfig
	Log    applog.LogConfig

	Mode string
}

type ServerConfig struct {
	ListenOn string
}

// Returns a new config and a setup viper instance if there is no error;
// If encounters an error, returns all nil except the error
func CreateAppConfig() (*AppConfig, *viper.Viper, error) {
	appConf := AppConfig{}
	vpr := viper.New()
	vpr.SetEnvPrefix("PORTFOLIO")

	if mode, err := GetMode(vpr); err != nil {
		return nil, nil, err
	} else {
		appConf.Mode = *mode
	}

	if serverConf, err := SetupServerConfig(vpr); err != nil {
		return nil, nil, err
	} else {
		appConf.Server = *serverConf
	}

	if dbConf, err := dbconfig.Create(vpr); err != nil {
		return nil, nil, err
	} else {
		appConf.DB = *dbConf
	}

	if logConf, err := applog.CreateLogConfig(vpr, appConf.Mode == "DEBUG"); err != nil {
		return nil, nil, err
	} else {
		appConf.Log = *logConf
	}

	return &appConf, vpr, nil
}

func GetMode(vpr *viper.Viper) (*string, error) {
	vpr.BindEnv("MODE")
	appmode := vpr.GetString("MODE")
	if appmode == "" {
		return nil, errors.New("Env variable MODE is unset")
	}

	return &appmode, nil
}

func SetupServerConfig(vpr *viper.Viper) (*ServerConfig, error) {
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

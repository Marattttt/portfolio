package config

import (
	"fmt"
	"strings"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/logconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/serverconfig"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config/storageconfig"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Storage storageconfig.StorageConfig
	Server  serverconfig.ServerConfig
	Log     logconfig.LogConfig

	Mode Mode
}

type Mode int

const (
	Debug Mode = iota
	Release
)

const (
	modeEnv = "MODE"
)

// Defaults
const (
	defMode = Debug
)

// Returns a new config and a setup viper instance if there is no error;
// If encounters an error, returns all nil except the error
func New() (*AppConfig, error) {
	conf := AppConfig{}
	vpr := viper.New()
	vpr.SetEnvPrefix("PORTFOLIO")

	mode, err := getMode(vpr)
	if err != nil {
		return nil, err
	}

	conf.Mode = mode

	server, err := serverconfig.New(vpr)
	if err != nil {
		return nil, fmt.Errorf("creating server config: %w", err)
	}
	conf.Server = *server

	db, err := storageconfig.New(vpr)
	if err != nil {
		return nil, fmt.Errorf("creating storage config: %w", err)
	}
	conf.Storage = *db

	isDebug := conf.Mode == Debug
	log, err := logconfig.New(vpr, isDebug)
	if err != nil {
		return nil, fmt.Errorf("creating log config: %w", err)
	}
	conf.Log = *log

	return &conf, nil
}

// Binds the key with viper.bindenv, returns the value in lowercase or nil if unset
func GetEnvString(vpr *viper.Viper, name string) *string {
	_ = vpr.BindEnv(name)
	res := strings.ToLower(vpr.GetString(name))
	if res == "" {
		return nil
	}
	return &res
}

func getMode(vpr *viper.Viper) (Mode, error) {
	m := GetEnvString(vpr, modeEnv)
	if m == nil {
		return defMode, nil
	}

	switch *m {
	case "debug":
		return Debug, nil
	case "release":
		return Release, nil
	default:
		return defMode, fmt.Errorf("Invalid value for env variable %s = %s", modeEnv, *m)
	}
}

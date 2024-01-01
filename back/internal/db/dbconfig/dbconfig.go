package dbconfig

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     uint
}

func (config DbConfig) GetDSN() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		config.Host,
		config.User,
		config.Password,
		config.DbName,
		config.Port)

	return dsn
}

// Uses following environment variables (prefix ommitted):
// "DB_HOSTNAME", "DB_PORT", "DB_USER", "DB_PASS", "DB_DBNAME"
func CreateConfig(vpr viper.Viper) (*DbConfig, error) {
	config := new(DbConfig)

	if err := fillConfig(config, vpr); err != nil {
		return nil, err
	}

	return config, nil
}

func fillConfig(conf *DbConfig, vpr viper.Viper) error {
	previousEnvPrefix := vpr.GetEnvPrefix()
	defer vpr.SetEnvPrefix(previousEnvPrefix)

	if previousEnvPrefix == "" {
		vpr.SetEnvPrefix("DB")
	} else {
		vpr.SetEnvPrefix(previousEnvPrefix + "_DB")
	}

	bindingErrorStart := "Error bindind env variable "
	if err := vpr.BindEnv("HOSTNAME"); err != nil {
		return errors.New(bindingErrorStart + "HOSTNAME")
	}
	if err := vpr.BindEnv("PORT"); err != nil {
		return errors.New(bindingErrorStart + "PORT")
	}
	if err := vpr.BindEnv("USER"); err != nil {
		return errors.New(bindingErrorStart + "USER")
	}
	if err := vpr.BindEnv("PASS"); err != nil {
		return errors.New(bindingErrorStart + "PASS")
	}
	if err := vpr.BindEnv("DBNAME"); err != nil {
		return errors.New(bindingErrorStart + "DBNAME")
	}

	conf.Host = vpr.GetString("HOSTNAME")
	conf.Port = vpr.GetUint("PORT")
	conf.User = vpr.GetString("USER")
	conf.Password = vpr.GetString("PASS")
	conf.DbName = vpr.GetString("DBNAME")

	return nil
}

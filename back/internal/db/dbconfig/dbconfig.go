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
	// Setting DB as env prefix is bad, because it messes up the viper key bindings

	bindingErrorStart := "Error bindind env variable "

	vpr.BindEnv("DB_HOSTNAME")
	vpr.BindEnv("DB_PORT")
	vpr.BindEnv("DB_USER")
	vpr.BindEnv("DB_PASS")
	vpr.BindEnv("DB_DBNAME")

	if conf.Host = vpr.GetString("DB_HOSTNAME"); conf.Host == "" {
		return errors.New(bindingErrorStart + conf.Host)
	}

	if conf.Port = vpr.GetUint("DB_PORT"); conf.Port == 0 {
		return errors.New(bindingErrorStart + fmt.Sprint(conf.Port))
	}

	if conf.User = vpr.GetString("DB_USER"); conf.User == "" {
		return errors.New(bindingErrorStart + conf.User)
	}

	if conf.Password = vpr.GetString("DB_PASS"); conf.Password == "" {
		return errors.New(bindingErrorStart + conf.Password)
	}

	if conf.DbName = vpr.GetString("DB_DBNAME"); conf.DbName == "" {
		return errors.New(bindingErrorStart + conf.DbName)
	}

	return nil
}

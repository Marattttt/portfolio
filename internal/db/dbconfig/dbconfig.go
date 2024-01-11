package dbconfig

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Stores a gorm.DB connection pointer, to share between requests, routines, etc.
type DbConfig struct {
	host     string
	user     string
	password string
	dbName   string
	port     uint

	dbconn *gorm.DB
}

// If the internal db connection is unset, opens a new one,
// otherwise, returns the existing connection
func (conf DbConfig) Connect() (*gorm.DB, error) {
	if conf.dbconn != nil {
		return conf.dbconn, nil
	}

	dsn := conf.getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		conf.dbconn = db
	}
	return conf.dbconn, err
}

func (config DbConfig) getDSN() string {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d",
		config.host,
		config.user,
		config.password,
		config.dbName,
		config.port)

	return dsn
}

// Uses following environment variables (prefix ommitted):
// "DB_HOSTNAME", "DB_PORT", "DB_USER", "DB_PASS", "DB_DBNAME"
func Create(vpr *viper.Viper) (*DbConfig, error) {
	config := new(DbConfig)

	if err := fillConfig(config, vpr); err != nil {
		return nil, err
	}

	return config, nil
}

func fillConfig(conf *DbConfig, vpr *viper.Viper) error {
	// Setting DB as env prefix is bad, because it messes up the viper key bindings

	bindingErrorStart := "Error bindind env variable "

	vpr.BindEnv("DB_HOSTNAME")
	vpr.BindEnv("DB_PORT")
	vpr.BindEnv("DB_USER")
	vpr.BindEnv("DB_PASS")
	vpr.BindEnv("DB_DBNAME")

	if conf.host = vpr.GetString("DB_HOSTNAME"); conf.host == "" {
		return errors.New(bindingErrorStart + conf.host)
	}

	if conf.port = vpr.GetUint("DB_PORT"); conf.port == 0 {
		return errors.New(bindingErrorStart + fmt.Sprint(conf.port))
	}

	if conf.user = vpr.GetString("DB_USER"); conf.user == "" {
		return errors.New(bindingErrorStart + conf.user)
	}

	if conf.password = vpr.GetString("DB_PASS"); conf.password == "" {
		return errors.New(bindingErrorStart + conf.password)
	}

	if conf.dbName = vpr.GetString("DB_DBNAME"); conf.dbName == "" {
		return errors.New(bindingErrorStart + conf.dbName)
	}

	return nil
}

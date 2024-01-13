package applog

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type LogConfig struct {
	OutPath     string
	IsDebugMode bool
}

type logger struct {
	Error *log.Logger
	Info  *log.Logger
}

var (
	applogger logger
	config    LogConfig
)

const (
	Config = "(Config) %v"
	Db     = "(DB) %v"
	Http   = "(HTTP) %v"
)

// Needs to be run before any use of the package
// Configures the logging behavior
func Setup(conf LogConfig) {
	config = conf

	logfile, err := os.OpenFile(conf.OutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error when opening file for logging\n", err)
	}

	applogger.Error = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	applogger.Info = log.New(logfile, "INFO: ", log.Ldate|log.Ltime)
}

func CreateLogConfig(vpr *viper.Viper, isDebugMode bool) (*LogConfig, error) {
	var newConfig LogConfig

	newConfig.IsDebugMode = isDebugMode

	vpr.BindEnv("LOGDEST")
	newConfig.OutPath = vpr.GetString("LOGDEST")
	if newConfig.OutPath == "" {
		return nil, fmt.Errorf("Env variable LOGDEST is unset")
	}

	return &newConfig, nil
}

func Error(format string, values ...interface{}) {
	applogger.Error.Printf(format, values...)

	if config.IsDebugMode {
		log.Printf(format, values...)
	}
}

func Info(format string, values ...interface{}) {
	applogger.Info.Printf(format, values...)

	if config.IsDebugMode {
		log.Printf(format, values...)
	}
}

func Fatal(format string, values ...interface{}) {
	if config.IsDebugMode {
		log.Printf(format, values...)
	}

	applogger.Error.Fatalf(format, values...)
}

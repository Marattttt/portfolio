package logconfig

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/configutils"
	"github.com/spf13/viper"
)

// All string values are lowercase
type LogConfig struct {
	LogDest     string
	Format      LogFormat
	IsDebugMode bool
	Flags       LogFlags
}

type LogFlags struct {
	Debug int
	Error int
	Info  int
	Warn  int
}

const (
	envPrefix = "LOG_"
	envFormat = envPrefix + "FORMAT"
	envDest   = envPrefix + "DEST"
)

// Defaults
const (
	defLogDest = "applog.log"
	defFormat  = JSONFormat

	defDebugFlags = log.Ltime | log.Lmicroseconds | log.Lshortfile
	defErrorFlags = log.Ldate | log.Ltime | log.Lshortfile
	defWarnFlags  = log.Ldate | log.Ltime | log.Lshortfile
	defInfoFlags  = log.Ldate | log.Ltime
)

func New(vpr *viper.Viper, isDebugMode bool) (*LogConfig, error) {
	var conf LogConfig
	conf.IsDebugMode = isDebugMode
	conf.setDefaultFlags()

	conf.setupDest(vpr)

	if err := conf.setupFormat(vpr); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *LogConfig) setupDest(vpr *viper.Viper) {
	if dest := configutils.GetEnvString(vpr, envDest); dest == nil {
		c.LogDest = defLogDest
	} else {
		c.LogDest = *dest
	}
}

func (c *LogConfig) setupFormat(vpr *viper.Viper) error {
	var formatStr *string

	allowed := map[string]LogFormat{
		"json": JSONFormat,
		"text": TextFormat,
	}

	if formatStr = configutils.GetEnvString(vpr, envFormat); formatStr == nil {
		c.Format = defFormat
		return nil
	}

	format, ok := allowed[*formatStr]

	if !ok {
		err := configutils.NewErrValueNotAllowed(envFormat, *formatStr, []string{"json", "text"})
		return err
	}

	c.Format = format

	return nil
}

func (c *LogConfig) setDefaultFlags() {
	c.Flags.Debug = defDebugFlags
	c.Flags.Error = defErrorFlags
	c.Flags.Info = defInfoFlags
	c.Flags.Warn = defWarnFlags
}

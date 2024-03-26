package logconfig

import (
	"log"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/configutils"
)

// All string values are lowercase
type LogConfig struct {
	Destination  string `env:"DEST, default=stdout"`
	FormatString string `env:"FORMAT, default=json"`
	Format       LogFormat
	IsDebugMode  bool
	Flags        LogFlags
}

type LogFormat int

const (
	JSONFormat = iota
	TextFormat
)

func (f LogFormat) String() string {
	switch f {
	case JSONFormat:
		return "json"
	case TextFormat:
		return "text"
	default:
		return "unidentified"
	}
}

func (f LogFormat) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}

type LogFlags struct {
	Debug int
	Error int
	Info  int
	Warn  int
}

// Log flags
const (
	debugFlags = log.Ltime | log.Lmicroseconds | log.Lshortfile
	errorFlags = log.Ldate | log.Ltime | log.Lshortfile
	warnFlags  = log.Ldate | log.Ltime | log.Lshortfile
	infoFlags  = log.Ldate | log.Ltime
)

func New(isDebugMode bool) (*LogConfig, error) {
	var conf LogConfig
	conf.IsDebugMode = isDebugMode
	conf.setDefaultFlags()

	if err := conf.setupFormat(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *LogConfig) setupFormat() error {
	allowed := map[string]LogFormat{
		"json": JSONFormat,
		"text": TextFormat,
	}

	if _, ok := allowed[c.FormatString]; !ok {
		err := configutils.NewErrValueNotAllowed(c.FormatString, c.FormatString, []string{"json", "text"})
		return err
	}

	c.Format = allowed[c.FormatString]

	return nil
}

func (c *LogConfig) setDefaultFlags() {
	c.Flags.Debug = debugFlags
	c.Flags.Error = errorFlags
	c.Flags.Info = infoFlags
	c.Flags.Warn = warnFlags
}

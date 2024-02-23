package serverconfig

import (
	"fmt"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/configutils"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ListenOn         string
	ReadTimout       time.Duration
	ReadHeaderTimout time.Duration
}

// Environment
const (
	envPort          = "PORT"
	envReadTimeot    = "REQ_READTIMEOUT"
	envHeaderTimeout = "REQ_READ_HEADER_TIMEOUT"
)

// Defaults
const (
	defPort          = 3000
	defReadTimout    = 30
	defHeaderTimeout = 3
)

func New(vpr *viper.Viper) (*ServerConfig, error) {
	var conf ServerConfig

	if p := configutils.GetEnvInt(vpr, envPort); p != nil {
		conf.ListenOn = fmt.Sprintf(":%d", *p)
	} else {
		conf.ListenOn = fmt.Sprintf(":%d", defPort)
	}

	if t := configutils.GetEnvInt(vpr, envReadTimeot); t != nil {
		conf.ReadTimout = time.Second * time.Duration(*t)
	} else {
		conf.ReadTimout = time.Second * defReadTimout
	}

	if ht := configutils.GetEnvInt(vpr, envHeaderTimeout); ht != nil {
		conf.ReadHeaderTimout = time.Second * time.Duration(*ht)
	} else {
		conf.ReadHeaderTimout = time.Second * defHeaderTimeout
	}

	return &conf, nil
}

package serverconfig

import (
	"fmt"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/configutils"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	ListenOn string
}

// Environment
const (
	envPort = "PORT"
)

// Defaults
const (
	defPort = 3000
)

func New(vpr *viper.Viper) (*ServerConfig, error) {
	var conf ServerConfig

	if p := configutils.GetEnvInt(vpr, envPort); p != nil {
		conf.ListenOn = fmt.Sprintf(":%d", *p)
	} else {
		conf.ListenOn = fmt.Sprintf(":%d", defPort)
	}

	return &conf, nil
}

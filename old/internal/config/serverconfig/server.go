package serverconfig

import (
	"time"
)

type ServerConfig struct {
	ListenOn          int           `env:"PORT, default=3030"`
	ReadTimout        time.Duration `env:"RTIMEOUT, default=20s"`
	ReadHeaderTimeout time.Duration `env:"RHEADERTIMEOUT, default=20s"`
}

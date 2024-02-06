package dbconfig

import "time"

type ConnParams struct {
	MaxConns          int
	MaxConnLifeTime   time.Duration
	HealthCheckPeriod time.Duration
}

type dbConnStr string

func (c dbConnStr) String() string {
	if isPGConnURL(string(c)) {
		return sanitizeConnURL(string(c))
	}

	return sanitizeDSN(string(c))
}

func (c dbConnStr) MarshalText() (text []byte, err error) {
	return []byte(c.String()), nil
}

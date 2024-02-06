package dbconfig

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/configutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// Stores a gorm.DB connection pointer, to share between requests, routines, etc.
type DbConfig struct {
	Pool     *pgxpool.Pool `json:"-"`
	poolConf *pgxpool.Config

	Connstr    dbConnStr
	ConnParams ConnParams
}

// Environment
const (
	envPrefix            = "DB_"
	envConnStr           = envPrefix + "CONN"
	envMaxConns          = envPrefix + "MAX_CONNS"
	envMaxConnLifeTime   = envPrefix + "MAX_CONN_LIFETIME"
	envHealthCheckPeriod = envPrefix + "HEALTH_CHECK_PERIOD"
)

// Defaults
const (
	defMaxConns          = 10
	defMaxConnLifetime   = time.Hour
	defHealthCheckPeriod = time.Minute
)

func New(vpr *viper.Viper) (*DbConfig, error) {
	var config DbConfig

	if err := config.fillParameters(vpr); err != nil {
		return nil, err
	}

	if err := config.addPool(); err != nil {
		return nil, err
	}

	return &config, nil
}

// Fills-out all external values
func (c *DbConfig) fillParameters(vpr *viper.Viper) error {
	connStr := configutils.GetEnvString(vpr, envConnStr)

	if connStr == nil {
		return fmt.Errorf("DSN not provided")
	}

	c.Connstr = dbConnStr(*connStr)

	c.fillConfigurations(vpr)

	return nil
}

func (c *DbConfig) fillConfigurations(vpr *viper.Viper) {
	params := &c.ConnParams
	if env := configutils.GetEnvInt(vpr, envMaxConns); env == nil {
		params.MaxConns = defMaxConns
	} else {
		params.MaxConns = *env
	}

	if env := configutils.GetEnvString(vpr, envMaxConnLifeTime); env == nil {
		params.MaxConnLifeTime = defMaxConnLifetime
	} else if parsed, err := time.ParseDuration(*env); err != nil {
		params.MaxConnLifeTime = defMaxConnLifetime
	} else {
		params.MaxConnLifeTime = parsed
	}

	if env := configutils.GetEnvString(vpr, envHealthCheckPeriod); env == nil {
		params.HealthCheckPeriod = defHealthCheckPeriod
	} else if parsed, err := time.ParseDuration(*env); err != nil {
		params.HealthCheckPeriod = defHealthCheckPeriod
	} else {
		params.HealthCheckPeriod = parsed
	}
}

// Should be called after configuring the poolConfig
func (c *DbConfig) addPool() error {
	if c.poolConf == nil {
		err := c.addPoolConfig()
		if err != nil {
			return err
		}
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), c.poolConf)
	if err != nil {
		return fmt.Errorf("while creating new pool with config: %w", err)
	}

	c.Pool = pool
	return nil
}

// Parses the poolconfig from dsn
func (c *DbConfig) addPoolConfig() error {
	c.setupDSN()

	dbconf, err := pgxpool.ParseConfig(string(c.Connstr))
	if err != nil {
		return fmt.Errorf("parsing dsn %s for pgxpool.Config: %w", c.Connstr, err)
	}

	c.poolConf = dbconf

	return nil
}

// DSN is setup with pgxpool additional parameters
func (c *DbConfig) setupDSN() {
	var format string

	//	postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	if isPGConnURL(string(c.Connstr)) {
		format = `?pool_max_conns=%d&pool_max_conn_lifetime=%s&pool_health_check_period=%s&`
	} else {
		format = " pool_max_conns=%d pool_max_conn_lifetime=%s pool_health_check_period=%s"
	}
	dsn := string(c.Connstr) + fmt.Sprintf(
		format,
		defMaxConns,
		defMaxConnLifetime,
		defHealthCheckPeriod,
	)

	c.Connstr = dbConnStr(dsn)
}

func isPGConnURL(connstr string) bool {
	r := regexp.MustCompile(`^[a-zA-Z]+://`)
	location := r.FindStringIndex(connstr)

	// Match is found and start is at the beginning of the string
	return location != nil && location[0] == 0
}

func sanitizeConnURL(url string) string {
	// From postgres://... match postgres://
	cleanrgx := regexp.MustCompile(`^.*\:\/\/`)

	// Trim the postgres:// part
	clean := cleanrgx.ReplaceAllString(url, "")

	// Isn't a connection string
	if clean == url {
		return url
	}

	sensitiveEnd := strings.IndexRune(clean, '@')
	passStart := strings.IndexRune(clean, ':')
	parametersStart := strings.IndexRune(clean, '/')

	if parametersStart == -1 {
		parametersStart = len(url) - 1
	}

	noSensitiveChars := sensitiveEnd == -1 || passStart == -1
	noSensitiveData := sensitiveEnd > parametersStart || passStart > parametersStart || passStart > sensitiveEnd

	if noSensitiveChars || noSensitiveData {
		return url
	}

	if sensitiveEnd > parametersStart && passStart > parametersStart {
		return url
	}

	userPass := clean[:sensitiveEnd]
	userNoPass := userPass[:strings.IndexRune(userPass, ':')]
	sanitizedUserPass := userNoPass + ":xxxxx"

	clean = strings.Replace(url, userPass, sanitizedUserPass, 1)

	return clean
}

func sanitizeDSN(dsn string) string {
	r := regexp.MustCompile(`[ ;]?[pP]assword=[^ ;]+`)

	return r.ReplaceAllString(dsn, "<password redacted>")
}

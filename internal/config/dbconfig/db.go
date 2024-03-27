package dbconfig

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Stores a gorm.DB connection pointer, to share between requests, routines, etc.
type DbConfig struct {
	Pool     *pgxpool.Pool `json:"-"`
	poolConf *pgxpool.Config

	// Url to use when using go-migrate
	// See https://github.com/golang-migrate/migrate
	MigrationsURL string `env:"MIGRATIONMS, default=file://migrations"`

	Conn              dbConnStr `env:"CONN, required"`
	ConnRaw           string    `json:"-" env:"CONN, required"`
	MaxConns          int       `env:"MAXCONNS, default=50"`
	MaxConnLifetime   int       `env:"MAXCONNLIFETIME, default=30"`
	HealthCheckPeriod int       `env:"HEALTHCHECKPERIOD, default=1"`
}

func (c *DbConfig) Close(_ context.Context) error {
	// FIXME: Cannot manually close a connecion on shutdown if it was not established successfully
	// if c.Pool != nil {
	// 	// Check for a working connection to detect possible issues
	// 	// and avoid possible errors/panics inside the lib's code if there are any problems
	// 	if err := c.Pool.Ping(ctx); err != nil {
	// 		return err
	// 	}
	// 	c.Pool.Close()
	// }
	return nil
}

func (c *DbConfig) Configure() error {
	var config DbConfig

	if err := config.setPool(); err != nil {
		return err
	}

	return nil
}

func (c *DbConfig) setPool() error {
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
	c.setupConnParams()

	dbconf, err := pgxpool.ParseConfig(string(c.Conn))
	if err != nil {
		return fmt.Errorf("parsing dsn %s for pgxpool.Config: %w", c.Conn, err)
	}

	c.poolConf = dbconf

	return nil
}

// DSN is setup with pgxpool additional parameters
func (c *DbConfig) setupConnParams() {
	var format = `?pool_max_conns=%d&pool_max_conn_lifetime=%d&pool_health_check_period=%d&`

	dsn := c.ConnRaw + fmt.Sprintf(
		format,
		c.MaxConns,
		c.MaxConnLifetime,
		c.HealthCheckPeriod,
	)

	c.Conn = dbConnStr(dsn)
}

func (c *DbConfig) Migrate(ctx context.Context, logger applog.Logger) error {
	m, err := migrate.New(
		c.MigrationsURL,
		string(c.Conn),
	)
	if err != nil {
		return fmt.Errorf("preparing for migration: %w", err)
	}

	logger.Info(ctx, applog.DB, "Beinning mirgation",
		slog.String("source", c.MigrationsURL),
		slog.String("connstr", c.ConnRaw))

	if err := m.Up(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

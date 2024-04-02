package applog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/Marattttt/portfolio/portfolio_back/internal/config/logconfig"
)

// All types implementing logger interface are safe to pass by reference
type Logger interface {
	Debug(ctx context.Context, s Scope, msg string, attrs ...slog.Attr)
	Error(ctx context.Context, s Scope, msg string, err error, attrs ...slog.Attr)
	Info(ctx context.Context, s Scope, msg string, attrs ...slog.Attr)
	Warn(ctx context.Context, s Scope, msg string, attrs ...slog.Attr)
	With(...any) Logger
}

type Scope int

const (
	//Scope flag for the application as a whole
	Application Scope = (1 << iota)
	// Scope flag for auth
	Auth
	// Scope flag for config
	Config
	// Scope flag for DB
	DB
	// Scope flag for logs not related to anythning
	Generic
	// Scope flag for HTTP
	HTTP

	scopeAttrKey = "scope"
)

// Parses scope flags
func (s Scope) LogValue() slog.Value {
	scopes := make([]string, 0)

	if s&Application == Application {
		scopes = append(scopes, "application")
	}
	if s&Config == Config {
		scopes = append(scopes, "config")
	}
	if s&DB == DB {
		scopes = append(scopes, "DB")
	}
	if s&HTTP == HTTP {
		scopes = append(scopes, "HTTP")
	}
	if s&Auth == Auth {
		scopes = append(scopes, "auth")
	}
	if s&Generic == Generic {
		scopes = append(scopes, "generic")
	}

	return slog.AnyValue(scopes)
}

func (s Scope) String() string {
	scopes := strings.Builder{}

	if s&Application == Application {
		scopes.WriteString("application")
	}
	if s&Config == Config {
		scopes.WriteString("config")
	}
	if s&DB == DB {
		scopes.WriteString("DB")
	}
	if s&HTTP == HTTP {
		scopes.WriteString("HTTP")
	}
	if s&Auth == Auth {
		scopes.WriteString("auth")
	}

	return scopes.String()
}

// Default logger impementation that wraps around the std slog package
type AppLogger struct {
	out  *slog.Logger
	conf *logconfig.LogConfig
}

// Example:
// logger.Error(ctx, applog.DB|applog.Config, "reading file", err, nil)
// or
// Example: logger.Error(ctx, applog.DB|applog.Config, "reading file", err, slog.String("somedata", data))
func (l AppLogger) Error(ctx context.Context, s Scope, msg string, err error, attrs ...slog.Attr) {
	scoped := l.out.With(scopeAttrKey, s)

	attrs = append(attrs, slog.Any("err", err))

	scoped.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

// Example: logger.Info(ctx, applog.DB|applog.Config, "reading file", slog.String("somedata", data))
func (l AppLogger) Info(ctx context.Context, s Scope, msg string, attrs ...slog.Attr) {
	scoped := l.out.With(scopeAttrKey, s)

	scoped.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

// Example: logger.Warn(ctx, applog.DB|applog.Config, "reading file", slog.String("somedata", data))
func (l AppLogger) Warn(ctx context.Context, s Scope, msg string, attrs ...slog.Attr) {
	scoped := l.out.With(scopeAttrKey, s)

	scoped.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
}

// Example: logger.Debug(ctx, applog.DB|applog.Config, "reading file", slog.String("somedata", data))
func (l AppLogger) Debug(ctx context.Context, s Scope, msg string, attrs ...slog.Attr) {
	if !l.conf.IsDebugMode {
		return
	}

	scoped := l.out.With(scopeAttrKey, s)

	scoped.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

// Example: requestLogger := logger.With("requestId", some-id)
func (l AppLogger) With(args ...any) Logger {
	l.out = l.out.With(args...)
	return l
}

// Create new logger
func New(conf logconfig.LogConfig) (*AppLogger, error) {
	var logger AppLogger
	logger.conf = &conf

	if err := logger.configureOutput(conf.Destination); err != nil {
		return nil, err
	}

	return &logger, nil
}

// If provided with "stdout" or "stderr" sets all outputs to that
// Any other string is treated as a path
func (l *AppLogger) configureOutput(to string) error {
	out, err := outputFor(to)
	if err != nil {
		return err
	}

	var h slog.Handler

	switch l.conf.Format {
	case logconfig.JSONFormat:
		h = slog.NewJSONHandler(*out, nil)
	case logconfig.TextFormat:
		h = slog.NewTextHandler(*out, nil)
	default:
		return fmt.Errorf("Invalid LogConfig.Format value %d", l.conf.Format)
	}

	l.out = slog.New(h)
	return nil
}

func outputFor(to string) (*io.Writer, error) {
	var out io.Writer

	switch to {
	case "stdout":
		out = os.Stdout
	case "stderr":
		out = os.Stderr
	default:
		o, err := getOutFile(to)
		if err != nil {
			return nil, err
		}
		out = o
	}

	return &out, nil
}

func getOutFile(name string) (*os.File, error) {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return nil, fmt.Errorf("getting output for %s: %w", name, err)
	}

	err = testOutFile(f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Tries writing 100 bytes to a file, checks for errors and truncates it to initial size
func testOutFile(f *os.File) error {
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("trying to call stat on log file: %w", err)
	}

	const writeLen = 100

	startSize := stat.Size()
	_, err = f.Write(make([]byte, writeLen))
	if err != nil {
		return fmt.Errorf("testing log file by writing %d bytes: %w", writeLen, err)
	}

	_ = f.Truncate(startSize)

	return nil
}

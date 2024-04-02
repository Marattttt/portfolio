package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	Adress int `env:"ADRESS, default=3030"`
	// Allowed values: "stderr", stdout or any other name to output to a file
	LogDestination string `env:"LOGDEST, default=stderr"`
	// Allowed values: "json", "text"
	LogFormat string `env:"LOGFORMAT, default=json"`
}

func createConfig(ctx context.Context) (*config, error) {
	conf := config{}
	if err := envconfig.Process(ctx, &conf); err != nil {
		return nil, fmt.Errorf("creting app config: %w", err)
	}
	return &conf, nil
}

func configureApp(conf config) error {

	output, err := configureLogOutput(conf.LogDestination)
	if err != nil {
		return err
	}

	handler, err := configureLogFormat(conf.LogFormat, output)
	if err != nil {
		return err
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

func configureLogOutput(destination string) (io.Writer, error) {
	var output io.Writer
	switch strings.ToLower(destination) {
	case "":
		output = os.Stderr
	case "stderr":
		output = os.Stderr
	case "stdout":
		output = os.Stdout
	default:
		fileout, err := os.OpenFile(destination, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return nil, fmt.Errorf("opening file for log output: %w", err)
		}
		output = fileout

	}
	return output, nil
}

func configureLogFormat(format string, output io.Writer) (slog.Handler, error) {
	var handler slog.Handler
	switch strings.ToLower(format) {
	case "text":
		handler = slog.NewTextHandler(output, nil)
	case "json":
		handler = slog.NewJSONHandler(output, nil)
	default:
		return nil, fmt.Errorf("unrecognized log format")
	}

	return handler, nil
}

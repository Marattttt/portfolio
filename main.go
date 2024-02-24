package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
)

func main() {
	cancelsignals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	appctx, appcancel := signal.NotifyContext(context.Background(), cancelsignals...)
	defer appcancel()

	conf := initAppConfig()
	printConfig(conf)

	logger := initLogger(&conf)

	server := api.NewServer(appctx, logger, &conf)

	// Server
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(appctx, applog.Application, "Unexpected server shutdown", err)
		}
	}()

	<-appctx.Done()

	const shutdownTimeout = time.Second * 2
	shutdownErrors := make(chan error, 100)
	shutdownctx, stopShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer stopShutdown()

	logger.Info(shutdownctx, applog.Application, "Beginning shutdown", slog.Duration("timeout", shutdownTimeout))

	go shutdownServer(shutdownctx, shutdownErrors, server)

	go shutdownConfig(shutdownctx, shutdownErrors, &conf)

	<-shutdownctx.Done()

	if shutdownctx.Err() == context.DeadlineExceeded {
		logger.Error(
			context.Background(),
			applog.Application,
			"Shutdown timed out",
			fmt.Errorf("shutting down took longer than %s", shutdownTimeout.String()))
	}

	printErrors(logger, shutdownErrors)
}

func initAppConfig() config.AppConfig {
	c, err := config.New()
	if err != nil {
		log.Fatalf("while creating appconfig: %v\n", err)
	}
	return *c
}

func initLogger(conf *config.AppConfig) applog.AppLogger {
	l, err := applog.New(conf.Log)
	if err != nil {
		log.Fatalf("while creating logger: %v", err)
	}
	return *l
}

func printConfig(conf config.AppConfig) {
	marshalledConf, err := json.MarshalIndent(conf, "", strings.Repeat(" ", 4))
	if err != nil {
		log.Fatalf("Marshalling created config: %v", err)
	}
	log.Println("Beginning start up using config: \n" + string(marshalledConf))
}

func shutdownServer(ctx context.Context, errs chan error, server *http.Server) {
	if err := server.Shutdown(ctx); err != nil {
		errs <- fmt.Errorf("shutting down http server: %w", err)
	}
}

func shutdownConfig(ctx context.Context, errs chan error, conf *config.AppConfig) {
	if err := conf.Close(ctx); err != nil {
		errs <- fmt.Errorf("closing config resources: %w", err)
	}
}

func printErrors(logger applog.Logger, errs chan error) {
	// Log all errrors
	errors := make([]error, 0, len(errs))
	for len(errs) > 0 {
		errors = append(errors, <-errs)
	}

	for _, err := range errors {
		logger.Error(context.Background(), applog.Application, "During shutdown", err)
	}
}

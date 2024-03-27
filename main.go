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
	"sync"
	"syscall"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
)

func main() {
	cancelsignals := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM}
	appCtx, appcancel := signal.NotifyContext(context.Background(), cancelsignals...)
	defer appcancel()

	conf := initAppConfig(appCtx)
	printConfig(conf)

	logger := initLogger(&conf)

	configure(appCtx, logger, conf)

	server := api.Server(appCtx, logger, &conf)

	go func() {
		serve(appCtx, logger, server)
		appcancel()
	}()

	// Wait for shutdown
	<-appCtx.Done()

	const shutdownTimeout = time.Second * 20

	shutdownErrors := make(chan error, 100)
	shutdownSuccess := make(chan struct{})

	shutdownWg := sync.WaitGroup{}

	timeoutCtx, stopShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer stopShutdown()

	logger.Info(timeoutCtx, applog.Application, "Beginning shutdown", slog.Duration("timeout", shutdownTimeout))

	shutdownWg.Add(1)
	go shutdownServer(timeoutCtx, &shutdownWg, shutdownErrors, server)
	shutdownWg.Add(1)
	go shutdownConfig(timeoutCtx, &shutdownWg, shutdownErrors, &conf)
	go func() {
		shutdownWg.Wait()
		shutdownSuccess <- struct{}{}
	}()

	waitShutdown(shutdownSuccess, timeoutCtx, logger)

	printErrors(logger, shutdownErrors)
}

func waitShutdown(success chan struct{}, timeout context.Context, logger applog.Logger) {
	select {
	case <-success:
		logger.Info(
			timeout,
			applog.Application,
			"Shutdown complete",
		)
	case <-timeout.Done():
		logger.Error(
			context.Background(),
			applog.Application,
			"Shutdown timed out",
			fmt.Errorf("Shutdown timeout exceeded"))
	}
}

func initAppConfig(ctx context.Context) config.AppConfig {
	c, err := config.New(ctx)
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

func serve(ctx context.Context, logger applog.Logger, server *http.Server) {
	logger.Info(ctx, applog.Application|applog.HTTP, "Starting server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(ctx, applog.Application, "Unexpected server shutdown", err)
	}
}

func configure(ctx context.Context, logger applog.Logger, conf config.AppConfig) {
	if err := conf.Configure(ctx, logger); err != nil {
		logger.Error(ctx, applog.Application, "Failed to configure on startup", err)
		os.Exit(1)
	}
}
func printConfig(conf config.AppConfig) {
	marshalledConf, err := json.MarshalIndent(conf, "", strings.Repeat(" ", 4))
	if err != nil {
		log.Fatalf("Marshalling created config: %v", err)
	}
	log.Println("Beginning start up using config: \n" + string(marshalledConf))
}

func shutdownServer(ctx context.Context, wg *sync.WaitGroup, errs chan error, server *http.Server) {
	defer wg.Done()
	if err := server.Shutdown(ctx); err != nil {
		errs <- fmt.Errorf("shutting down http server: %w", err)
	}
}

func shutdownConfig(ctx context.Context, wg *sync.WaitGroup, errs chan error, conf *config.AppConfig) {
	defer wg.Done()
	if err := conf.Close(ctx); err != nil {
		errs <- fmt.Errorf("closing config resources: %w", err)
	}
}

func printErrors(logger applog.Logger, errs chan error) {
	errcount := len(errs)
	for i := 0; i < errcount; i++ {
		err := <-errs
		logger.Error(context.Background(), applog.Application, "During shutdown", err)
	}
}

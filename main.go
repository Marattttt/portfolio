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

	server := api.Server(appCtx, logger, &conf)

	// Serve
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(appCtx, applog.Application, "Unexpected server shutdown", err)
		}
		appcancel()
	}()

	// Debug
	go func() {
		counter := 0
		for {
			select {
			case <-appCtx.Done():
				return
			case <-time.After(time.Second * 2):
				logger.Debug(context.Background(), applog.Application, "Still working", slog.Int("counter", counter))
				counter++
			}
		}
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

	select {
	case <-shutdownSuccess:
		logger.Info(
			timeoutCtx,
			applog.Application,
			"Shutdown complete",
		)
	case <-timeoutCtx.Done():
		logger.Error(
			context.Background(),
			applog.Application,
			"Shutdown timed out",
			fmt.Errorf("shutting down took longer than %s", shutdownTimeout.String()))
	}

	printErrors(logger, shutdownErrors)
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

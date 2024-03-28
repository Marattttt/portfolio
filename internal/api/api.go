package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/apigen"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type apiServerCodegenWrapper struct {
	apigen.Unimplemented
}

type requestData struct {
	id     uint64
	logger applog.Logger
}

type statusCodeResponseWriter struct {
	http.ResponseWriter
}

var (
	conf   *config.AppConfig
	served atomic.Uint64
)

func Server(basectx context.Context, logger applog.Logger, appconfig *config.AppConfig) *http.Server {
	conf = appconfig

	mux := chi.NewMux()
	mux.Use(middleware.Recoverer)

	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			data := ctx.Value(requestData{}).(requestData)

			data.logger.Info(ctx, applog.HTTP, "New request")

			start := time.Now()

			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrapped, r)

			took := time.Since(start)
			data.logger.Info(ctx,
				applog.HTTP,
				"Finished request",
				slog.Duration("timeTook", took),
				slog.Int("responseCode", wrapped.Status()))
		})
	})

	handler := apigen.HandlerFromMux(apiServerCodegenWrapper{}, mux)

	address := fmt.Sprintf(":%d", conf.Server.ListenOn)

	server := http.Server{
		Addr:              address,
		ReadHeaderTimeout: conf.Server.ReadHeaderTimeout,
		ReadTimeout:       conf.Server.ReadTimout,
		Handler:           handler,

		BaseContext: func(_ net.Listener) context.Context {
			requestID := served.Add(1)
			logCtx := context.WithValue(basectx, requestData{}, requestData{
				id:     requestID,
				logger: logger.With(slog.Uint64("requestID", requestID)),
			})
			return logCtx
		},
	}

	return &server
}

// func (apiServerCodegenWrapper) PostRegistser(w http.ResponseWriter, r *http.Request) {
// 	const maxBodyLength = 4000
// 	var (
// 		ctx          = r.Context()
// 		logger       = ctx.Value("logger").(applog.Logger)
// 		registerData = &apigen.RegisterRequest{}
// 		body         = make([]byte, 0)
// 	)
// 	if _, err := r.Body.Read(body); err != nil {
// 		logger.Error(ctx, applog.HTTP, "Failed to read request body", err)
// 	}
// }

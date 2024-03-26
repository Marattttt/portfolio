package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/apigen"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/go-chi/chi/v5"
)

type apiServerCodegenWrapper struct {
	apigen.Unimplemented
}

type requestData struct {
	id     uint64
	logger applog.Logger
}

var (
	conf   *config.AppConfig
	served atomic.Uint64
)

func Server(basectx context.Context, logger applog.Logger, appconfig *config.AppConfig) *http.Server {
	conf = appconfig

	mux := chi.NewMux()
	// mux.Use(logServedRequest)

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

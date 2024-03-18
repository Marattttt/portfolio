package api

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/Marattttt/portfolio/portfolio_back/internal/guests"
	"github.com/go-chi/chi/v5"
)

type apiServerCodegenWrapper struct {
}

var (
	conf   *config.AppConfig
	logger applog.Logger
	served atomic.Uint64
)

func NewMux(basectx context.Context, l applog.Logger, c *config.AppConfig) http.Handler {
	logger = l
	conf = c

	mux := chi.NewMux()
	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info(r.Context(), applog.HTTP, "Beginning serving request")
			h.ServeHTTP(w, r)
			l.Info(r.Context(), applog.HTTP, "Finished serving request")
		})
	})

	// handler := apigen.HandlerFromMux(apiServerCodegenWrapper{}, mux)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healthy!"))
	})

	_ = http.Server{
		Addr: conf.Server.ListenOn,
		// ReadHeaderTimeout: conf.Server.ReadHeaderTimout,
		// ReadTimeout:       conf.Server.ReadTimout,
		// Handler:           handler,
		Handler: mux,
		// BaseContext: func(_ net.Listener) context.Context {
		// 	return basectx
		// },
	}

	return mux
}

func (apiServerCodegenWrapper) GetGuestsGuestId(w http.ResponseWriter, r *http.Request, guestId int) {
	ctx := r.Context()

	guests, err := guests.NewFromConfig(logger, conf)
	if err != nil {
		logger.Error(ctx, applog.Application|applog.Config, "failed to create guests service from config", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	reuested := guests.GetGuest(ctx, guestId)
	if reuested == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(reuested)
	if err != nil {
		logger.Error(ctx, applog.Generic, "Failed to encode guest model", err)
	}

	_, _ = w.Write(data)
}

func (apiServerCodegenWrapper) PostAuthorize(w http.ResponseWriter, r *http.Request) {

}
func (apiServerCodegenWrapper) PostGuests(w http.ResponseWriter, r *http.Request) {}

func (apiServerCodegenWrapper) PatchGuestsGuestId(w http.ResponseWriter, r *http.Request, guestId int) {
}

package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/apigen"
	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/Marattttt/portfolio/portfolio_back/internal/config"
	"github.com/Marattttt/portfolio/portfolio_back/internal/users"
	"github.com/go-chi/chi/v5"
)

type apiServerCodegenWrapper struct {
	apigen.Unimplemented
}

var (
	conf   *config.AppConfig
	logger applog.Logger
	served atomic.Uint64
)

func Server(basectx context.Context, l applog.Logger, c *config.AppConfig) *http.Server {
	logger = l
	conf = c

	mux := chi.NewMux()
	mux.Use(logServedRequest)

	handler := apigen.HandlerFromMux(apiServerCodegenWrapper{}, mux)

	server := http.Server{
		Addr:              conf.Server.ListenOn,
		ReadHeaderTimeout: conf.Server.ReadHeaderTimout,
		ReadTimeout:       conf.Server.ReadTimout,
		Handler:           handler,
		BaseContext: func(_ net.Listener) context.Context {
			return basectx
		},
	}

	return &server
}

func (apiServerCodegenWrapper) GetGuestsGuestId(w http.ResponseWriter, r *http.Request, guestId int) {
	ctx := r.Context()

	guests, err := users.NewFromConfig(logger, conf)
	if err != nil {
		logger.Error(ctx, applog.Application|applog.Config, "failed to create guests service from config", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	reuested := guests.Get(ctx, guestId)
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

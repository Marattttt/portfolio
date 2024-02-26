package api

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/Marattttt/portfolio/portfolio_back/internal/api/apigen"
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
)

func NewServer(basectx context.Context, l applog.Logger, c *config.AppConfig) *http.Server {
	logger = l
	conf = c

	mux := chi.NewMux()
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

	guests, err := guests.NewFromConfig(logger, conf)
	if err != nil {
		logger.Error(ctx, applog.Application|applog.Config, "failed to create guests service from config", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	reqquested := guests.GetGuest(ctx, guestId)
	if reqquested == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	data, err := json.Marshal(reqquested)
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

package api

import (
	"log/slog"
	"net/http"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
)

func logServedRequest(h http.Handler, logger applog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := served.Add(1)

		logger.Info(r.Context(), applog.HTTP, "Beginning serving request", slog.Int64("id", int64(id)))
		h.ServeHTTP(w, r)
		logger.Info(r.Context(), applog.HTTP, "Finished serving request", slog.Int64("id", int64(id)))
	})
}

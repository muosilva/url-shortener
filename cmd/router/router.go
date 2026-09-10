package router

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/muosilva/url-shortener/internal/handler"
	"github.com/muosilva/url-shortener/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(metrics.HTTPMiddleware)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("health check requested", "method", r.Method, "path", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	r.Get("/{code}", h.RedirectToOriginalURL)

	r.Route("/urls", func(r chi.Router) {
		r.Post("/", h.CreateURL)
	})

	r.Handle("/metrics", promhttp.Handler())

	return r
}

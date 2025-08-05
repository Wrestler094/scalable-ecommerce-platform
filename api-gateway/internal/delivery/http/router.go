package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	ProxyHandler      ProxyHandler
	MonitoringHandler *MonitoringHandler
}

func NewRouter(h Handlers, routes []string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		for _, path := range routes {
			r.Mount("/"+path, h.ProxyHandler.HandlerFor(path))
		}
	})

	// Infra namespace (Monitoring endpoints)
	r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	r.Get("/healthz", h.MonitoringHandler.Liveness)
	r.Get("/readyz", h.MonitoringHandler.Readiness)

	return r
}

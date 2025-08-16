package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/delivery/http/infra"
	"github.com/Wrestler094/scalable-ecommerce-platform/payment-service/internal/delivery/http/v1"
)

type Handlers struct {
	V1Handlers        v1.Handlers
	MonitoringHandler *infra.MonitoringHandler
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API namespace
	r.Route("/api", func(r chi.Router) {
		// v1 namespace
		r.Mount("/v1", v1.NewV1Router(h.V1Handlers))
	})

	// Infra namespace (Monitoring endpoints)
	r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	r.Get("/healthz", h.MonitoringHandler.Liveness)
	r.Get("/readyz", h.MonitoringHandler.Readiness)

	return r
}

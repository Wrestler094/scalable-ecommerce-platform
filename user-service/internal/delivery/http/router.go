package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	UserHandler       *UserHandler
	MonitoringHandler *MonitoringHandler
}

func NewRouter(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// API namespace
	r.Route("/api", func(r chi.Router) {
		// User auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.UserHandler.Register)
			r.Post("/login", h.UserHandler.Login)
			r.Post("/refresh", h.UserHandler.Refresh)
			r.Post("/logout", h.UserHandler.Logout)
		})
	})

	// Monitoring endpoints
	r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	r.Get("/healthz", h.MonitoringHandler.Liveness)
	r.Get("/readyz", h.MonitoringHandler.Readiness)

	return r
}

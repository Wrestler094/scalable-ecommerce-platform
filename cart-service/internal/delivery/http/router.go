package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	CartHandler       *CartHandler
	MonitoringHandler *MonitoringHandler
}

func NewRouter(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// API namespace
	r.Route("/api", func(r chi.Router) {
		// Cart routes
		r.Route("/cart", func(r chi.Router) {
			// Authorized only
			r.Group(authorizedOnly(authenticatorImpl, func(r chi.Router) {
				r.Get("/", h.CartHandler.GetCart)
				r.Post("/", h.CartHandler.AddItem)
				r.Put("/", h.CartHandler.UpdateItem)
				r.Delete("/", h.CartHandler.RemoveItem)
				r.Delete("/clear", h.CartHandler.ClearCart)
			}))
		})
	})

	// Monitoring endpoints
	r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	r.Get("/healthz", h.MonitoringHandler.Liveness)
	r.Get("/readyz", h.MonitoringHandler.Readiness)

	return r
}

func authorizedOnly(auth authenticator.Authenticator, handler func(r chi.Router)) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authenticator.RequireRoles(auth, authenticator.User, authenticator.Admin))
		handler(r)
	}
}

package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	ProductHandler    *ProductHandler
	CategoryHandler   *CategoryHandler
	MonitoringHandler *MonitoringHandler
}

func NewRouter(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// API namespace
	r.Route("/api", func(r chi.Router) {
		// Product routes
		r.Route("/products", func(r chi.Router) {
			// Public
			r.Get("/{id}", h.ProductHandler.GetProductByID)

			// Authorized only
			r.Group(adminOnly(authenticatorImpl, func(r chi.Router) {
				r.Post("/", h.ProductHandler.CreateProduct)
				r.Put("/{id}", h.ProductHandler.UpdateProduct)
				r.Delete("/{id}", h.ProductHandler.DeleteProduct)
			}))
		})

		// Category routes
		r.Route("/categories", func(r chi.Router) {
			// Public
			r.Get("/", h.CategoryHandler.GetAllCategories)
			r.Get("/{id}/products", h.CategoryHandler.GetProductsByCategoryID)

			// Admin only
			r.Group(adminOnly(authenticatorImpl, func(r chi.Router) {
				r.Post("/", h.CategoryHandler.CreateCategory)
			}))
		})
	})

	// Monitoring endpoints
	r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	r.Get("/healthz", h.MonitoringHandler.Liveness)
	r.Get("/readyz", h.MonitoringHandler.Readiness)

	return r
}

func adminOnly(auth authenticator.Authenticator, handler func(r chi.Router)) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authenticator.RequireRoles(auth, authenticator.Admin))
		handler(r)
	}
}

func authorizedOnly(auth authenticator.Authenticator, handler func(r chi.Router)) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authenticator.RequireRoles(auth, authenticator.User, authenticator.Admin))
		handler(r)
	}
}

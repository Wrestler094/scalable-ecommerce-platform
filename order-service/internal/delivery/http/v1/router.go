package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	OrderHandler *OrderHandler
}

func NewV1Router(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Orders routes
	r.Route("/orders", func(r chi.Router) {
		r.Use(authenticator.RequireAuth())

		r.Get("/", h.OrderHandler.GetOrdersList)
		r.Post("/", h.OrderHandler.CreateOrder)
		r.Get("/{id}", h.OrderHandler.GetOrderByID)
	})

	return r
}

package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	OrderHandler *OrderHandler
}

func NewV1Router(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

	// Orders routes
	r.Route("/orders", func(r chi.Router) {
		r.Group(authorizedOnly(authenticatorImpl, func(r chi.Router) {
			r.Get("/", h.OrderHandler.GetOrdersList)
			r.Post("/", h.OrderHandler.CreateOrder)
			r.Get("/{id}", h.OrderHandler.GetOrderByID)
		}))
	})

	return r
}

func authorizedOnly(auth authenticator.Authenticator, handler func(r chi.Router)) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authenticator.RequireRoles(auth, authenticator.User, authenticator.Admin))
		handler(r)
	}
}
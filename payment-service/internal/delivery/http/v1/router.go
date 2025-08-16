package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	PaymentHandler *PaymentHandler
}

func NewV1Router(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Payment routes
	r.Route("/payments", func(r chi.Router) {
		r.Use(authenticator.RequireAuth())

		r.Post("/pay", h.PaymentHandler.Pay)
	})

	return r
}

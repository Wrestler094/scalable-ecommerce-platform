package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	PaymentHandler *PaymentHandler
}

func NewV1Router(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

	// Payment routes
	r.Route("/payments", func(r chi.Router) {
		// Authorized only
		r.Group(authorizedOnly(authenticatorImpl, func(r chi.Router) {
			r.Post("/pay", h.PaymentHandler.Pay)
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
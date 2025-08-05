package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	CartHandler *CartHandler
}

func NewV1Router(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

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

	return r
}

func authorizedOnly(auth authenticator.Authenticator, handler func(r chi.Router)) func(r chi.Router) {
	return func(r chi.Router) {
		r.Use(authenticator.RequireRoles(auth, authenticator.User, authenticator.Admin))
		handler(r)
	}
}
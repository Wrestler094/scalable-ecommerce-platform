package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	UserHandler *UserHandler
}

func NewV1Router(h Handlers) http.Handler {
	r := chi.NewRouter()

	// User auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.UserHandler.Register)
		r.Post("/login", h.UserHandler.Login)
		r.Post("/refresh", h.UserHandler.Refresh)
		r.Post("/logout", h.UserHandler.Logout)
	})

	return r
}

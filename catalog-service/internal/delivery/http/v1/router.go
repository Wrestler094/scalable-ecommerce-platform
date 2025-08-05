package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
)

type Handlers struct {
	ProductHandler  *ProductHandler
	CategoryHandler *CategoryHandler
}

func NewV1Router(h Handlers, authenticatorImpl authenticator.Authenticator) http.Handler {
	r := chi.NewRouter()

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
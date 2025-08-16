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

func NewV1Router(h Handlers) http.Handler {
	r := chi.NewRouter()

	// Product routes
	r.Route("/products", func(r chi.Router) {
		// Public endpoints
		r.Get("/{id}", h.ProductHandler.GetProductByID)

		// Admin only endpoints
		r.Group(func(r chi.Router) {
			r.Use(authenticator.RequireAdmin())

			r.Post("/", h.ProductHandler.CreateProduct)
			r.Put("/{id}", h.ProductHandler.UpdateProduct)
			r.Delete("/{id}", h.ProductHandler.DeleteProduct)
		})
	})

	// Category routes
	r.Route("/categories", func(r chi.Router) {
		// Public endpoints
		r.Get("/", h.CategoryHandler.GetAllCategories)
		r.Get("/{id}/products", h.CategoryHandler.GetProductsByCategoryID)

		// Admin only endpoints
		r.Group(func(r chi.Router) {
			r.Use(authenticator.RequireAdmin())

			r.Post("/", h.CategoryHandler.CreateCategory)
		})
	})

	return r
}

package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	// No handlers yet, but keeping structure for future expansion
}

func NewV1Router(h Handlers) http.Handler {
	r := chi.NewRouter()

	// No API endpoints yet, but ready for future business logic
	// When adding notification endpoints, they would go here

	return r
}
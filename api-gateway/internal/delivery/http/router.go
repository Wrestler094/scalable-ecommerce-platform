package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(proxyHandler ProxyHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/user", proxyHandler.HandlerFor("user"))
		r.Mount("/order", proxyHandler.HandlerFor("order"))
		r.Mount("/catalog", proxyHandler.HandlerFor("catalog"))
		r.Mount("/cart", proxyHandler.HandlerFor("cart"))
		r.Mount("/payment", proxyHandler.HandlerFor("payment"))
		r.Mount("/notification", proxyHandler.HandlerFor("notification"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}

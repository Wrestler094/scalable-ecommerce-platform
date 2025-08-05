package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(proxyHandler ProxyHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/user", proxyHandler.HandlerFor("user"))
		r.Mount("/order", proxyHandler.HandlerFor("order"))
		r.Mount("/catalog", proxyHandler.HandlerFor("catalog"))
		r.Mount("/cart", proxyHandler.HandlerFor("cart"))
		r.Mount("/payment", proxyHandler.HandlerFor("payment"))
		r.Mount("/notification", proxyHandler.HandlerFor("notification"))
	})

	// Infra namespace (Monitoring endpoints)
	//r.Handle("/metrics", http.HandlerFunc(h.MonitoringHandler.Metrics))
	//r.Get("/healthz", h.MonitoringHandler.Liveness)
	//r.Get("/readyz", h.MonitoringHandler.Readiness)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return r
}

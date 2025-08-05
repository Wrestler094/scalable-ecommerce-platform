package infra

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/healthcheck"
)

type MonitoringHandler struct {
	manager healthcheck.Manager
}

func NewMonitoringHandler(manager healthcheck.Manager) *MonitoringHandler {
	return &MonitoringHandler{manager: manager}
}

func (h *MonitoringHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
}

func (h *MonitoringHandler) Liveness(w http.ResponseWriter, _ *http.Request) {
	if !h.manager.IsAlive() {
		http.Error(w, "service not alive", http.StatusServiceUnavailable)
		return
	}

	respondText(w, http.StatusOK, "ok")
}

func (h *MonitoringHandler) Readiness(w http.ResponseWriter, _ *http.Request) {
	if !h.manager.IsReady() {
		http.Error(w, "service not ready", http.StatusServiceUnavailable)
		return
	}

	respondText(w, http.StatusOK, "ok")
}

func respondText(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

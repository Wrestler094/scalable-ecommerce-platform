package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
)

type ProxyHandler interface {
	HandlerFor(service string) http.Handler
}

type StaticProxyHandler struct {
	Targets map[string]string
	Logger  logger.Logger
}

func NewStaticProxyHandler(targets map[string]string, logger logger.Logger) *StaticProxyHandler {
	return &StaticProxyHandler{
		Targets: targets,
		Logger:  logger,
	}
}

func (h *StaticProxyHandler) HandlerFor(service string) http.Handler {
	const op = "StaticProxyHandler.HandlerFor"
	log := h.Logger.WithOp(op).With("service", service)

	targetURLStr, ok := h.Targets[service]
	if !ok {
		log.Warn("unknown service")
		return http.NotFoundHandler()
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.Warn("Invalid target URL", "url", targetURLStr, "error", err)
		return http.NotFoundHandler()
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		prefix := "/api/" + service
		trimmed := strings.TrimPrefix(req.URL.Path, prefix)

		if trimmed == "" || trimmed[0] != '/' {
			trimmed = "/" + trimmed
		}

		req.URL.Path = "/api" + trimmed
	}

	return proxy
}

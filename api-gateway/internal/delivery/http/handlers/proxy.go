package handlers

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

type ProxyHandler struct {
	Targets map[string]string
	Logger  logger.Logger
}

func NewProxyHandler(targets map[string]string, logger logger.Logger) *ProxyHandler {
	return &ProxyHandler{
		Targets: targets,
		Logger:  logger,
	}
}

func (h *ProxyHandler) HandlerFor(service string) http.Handler {
	const op = "ProxyHandler.HandlerFor"
	log := h.Logger.WithOp(op).With("service", service)

	targetURLStr, ok := h.Targets[service]
	if !ok {
		log.Warn("unknown service")
		return http.NotFoundHandler()
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.WithError(err).Warn("Invalid target URL", "url", targetURLStr)
		return http.NotFoundHandler()
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// TODO: GET FROM CONFIG AND IMPROVE
	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		if requestID := middleware.GetReqID(req.Context()); requestID != "" {
			req.Header.Set("X-Request-ID", requestID)
		}

		prefix := "/api/" + service
		trimmed := strings.TrimPrefix(req.URL.Path, prefix)

		if trimmed == "" || trimmed[0] != '/' {
			trimmed = "/" + trimmed
		}

		req.URL.Path = "/api" + trimmed
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		requestID := middleware.GetReqID(r.Context())
		log.WithError(err).
			WithRequestID(requestID).
			Error("Proxy error", "method", r.Method, "path", r.URL.Path)

		switch {
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Gateway Timeout", http.StatusGatewayTimeout)
		case errors.Is(err, context.Canceled):
			// Клиент отменил запрос
			return
		case strings.Contains(err.Error(), "connection refused"),
			strings.Contains(err.Error(), "no such host"):
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		default:
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}
	}

	return proxy
}

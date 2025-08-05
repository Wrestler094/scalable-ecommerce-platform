package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler interface {
	HandlerFor(service string) http.Handler
}

type StaticProxyHandler struct {
	Targets map[string]string
}

func NewStaticProxyHandler(targets map[string]string) *StaticProxyHandler {
	return &StaticProxyHandler{Targets: targets}
}

func (h *StaticProxyHandler) HandlerFor(service string) http.Handler {
	targetURLStr, ok := h.Targets[service]
	if !ok {
		log.Fatalf("unknown service: %s", service)
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.Fatalf("invalid URL for %s: %v", service, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
	}

	return proxy
}

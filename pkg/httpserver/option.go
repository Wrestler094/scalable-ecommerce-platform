package httpserver

import (
	"net/http"
	"time"
)

// Option represents a functional option for configuring the HTTP server.
type Option func(*Server)

// Port sets the TCP address (host:port) the server will listen on.
func Port(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

// Handler sets the HTTP handler used to process incoming requests.
// This option is required; otherwise, Start() will return an error.
func Handler(h http.Handler) Option {
	return func(s *Server) {
		s.App = h
	}
}

// ReadTimeout sets the maximum duration for reading the entire request, including the body.
func ReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = t
	}
}

// WriteTimeout sets the maximum duration before timing out writes of the response.
func WriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = t
	}
}

// ShutdownTimeout sets the timeout for gracefully shutting down the server.
func ShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}

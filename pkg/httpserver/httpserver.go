package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

const (
	_defaultAddr            = ":80"
	_defaultReadTimeout     = 5 * time.Second
	_defaultWriteTimeout    = 5 * time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

// Server wraps an HTTP server with configuration, start/stop logic, and error notifications.
type Server struct {
	// App is the HTTP handler to process requests. Must not be nil.
	App    http.Handler
	notify chan error

	address         string
	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	server *http.Server
}

// NewServer creates a new HTTP server with the given options.
func NewServer(opts ...Option) *Server {
	s := &Server{
		address:         _defaultAddr,
		readTimeout:     _defaultReadTimeout,
		writeTimeout:    _defaultWriteTimeout,
		shutdownTimeout: _defaultShutdownTimeout,
		notify:          make(chan error, 1),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start launches the HTTP server in a separate goroutine.
// Returns an error if the handler is not set or if initial configuration fails.
func (s *Server) Start() error {
	if s.App == nil {
		return errors.New("httpserver: handler is nil – use the Handler(...) option to provide a valid http.Handler")
	}

	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.App,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}

	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()

	return nil
}

// Notify returns a channel that receives any errors from ListenAndServe.
// This may include http.ErrServerClosed if shutdown was triggered manually. Caller is responsible for filtering if needed.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown gracefully stops the HTTP server using the configured shutdown timeout.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

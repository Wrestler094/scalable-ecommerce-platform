package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	_defaultAddr            = ":50051"
	_defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	App             *grpc.Server
	notify          chan error
	address         string
	shutdownTimeout time.Duration
	grpcOpts        []grpc.ServerOption
}

func New(opts ...Option) *Server {
	s := &Server{
		notify:          make(chan error, 1),
		address:         _defaultAddr,
		shutdownTimeout: _defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.App = grpc.NewServer(s.grpcOpts...)

	return s
}

func (s *Server) Start() {
	go func() {
		listener, err := net.Listen("tcp", s.address)
		if err != nil {
			s.notify <- fmt.Errorf("failed to listen on %s: %w", s.address, err)
			close(s.notify)
			return
		}

		s.notify <- s.App.Serve(listener)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.App.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.App.Stop()
		return fmt.Errorf("grpc server shutdown timed out")
	case <-stopped:
		return nil
	}
}

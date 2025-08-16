package grpcserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.address = net.JoinHostPort("", port)
	}
}

// ShutdownTimeout устанавливает время ожидания для graceful shutdown.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}

// UnaryInterceptor добавляет унарный перехватчик для всех RPC.
// Можно вызывать несколько раз для добавления нескольких перехватчиков.
func UnaryInterceptor(interceptor grpc.UnaryServerInterceptor) Option {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, grpc.UnaryInterceptor(interceptor))
	}
}

// StreamInterceptor добавляет стриминговый перехватчик.
func StreamInterceptor(interceptor grpc.StreamServerInterceptor) Option {
	return func(s *Server) {
		s.grpcOpts = append(s.grpcOpts, grpc.StreamInterceptor(interceptor))
	}
}

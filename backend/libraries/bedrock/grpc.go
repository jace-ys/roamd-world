package bedrock

import (
	"context"
	"fmt"
	"net"

	"github.com/alexliesenfeld/health"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/providers/zap/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jace-ys/roamd-world/backend/libraries/healthcheck"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type GRPCServer struct {
	srv          *grpc.Server
	healthchecks []health.Check

	addr string
}

func NewGRPCServer(logger *zap.SugaredLogger, port int) *GRPCServer {
	addr := fmt.Sprintf(":%d", port)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			logging.UnaryServerInterceptor(
				grpczap.InterceptorLogger(logger.Desugar()),
				logging.WithLevels(func(code codes.Code) logging.Level {
					return logging.DEBUG
				}),
			),
		),
	}

	return &GRPCServer{
		srv: grpc.NewServer(opts...),
		healthchecks: []health.Check{
			healthcheck.GRPCCheck("bedrock", addr),
		},
		addr: addr,
	}
}

func (s *GRPCServer) Name() string {
	return "grpc"
}

func (s *GRPCServer) Addr() string {
	return s.addr
}

func (s *GRPCServer) HealthChecks() []health.Check {
	return s.healthchecks
}

type ServiceHandler interface {
	HealthChecks() []health.Check
}

func (s *GRPCServer) RegisterService(sd *grpc.ServiceDesc, h ServiceHandler) {
	s.srv.RegisterService(sd, h)
	s.healthchecks = append(s.healthchecks, h.HealthChecks()...)
}

func (s *GRPCServer) Serve(ctx context.Context) error {
	healthpb.RegisterHealthServer(s.srv, &HealthHandler{})

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc server failed to serve: %w", err)
	}

	return s.srv.Serve(lis)
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	ok := make(chan struct{})

	go func() {
		s.srv.GracefulStop()
		close(ok)
	}()

	select {
	case <-ok:
		return nil
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	}
}

type HealthHandler struct {
}

func (h *HealthHandler) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: healthpb.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthHandler) Watch(req *healthpb.HealthCheckRequest, server healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

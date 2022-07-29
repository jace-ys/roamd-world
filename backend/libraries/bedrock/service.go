package bedrock

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	log     *zap.SugaredLogger
	servers []Server
}

func NewService(logger *zap.SugaredLogger, servers ...Server) *Service {
	return &Service{
		log:     logger,
		servers: append([]Server{}, servers...),
	}
}

type Server interface {
	Name() string
	Addr() string
	Serve(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func (s *Service) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			s.log.Infow("server listening", "name", srv.Name(), "addr", srv.Addr())
			return srv.Serve(ctx)
		})
	}

	s.log.Infow("service started")
	<-ctx.Done()
	stop()
	s.log.Infow("service shutting down gracefully")

	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, srv := range s.servers {
		srv := srv
		g.Go(func() error {
			if err := srv.Shutdown(sctx); err != nil {
				s.log.Warnw("server shutdown error", "name", srv.Name(), "error", err)
			} else {
				s.log.Infow("server shutdown complete", "name", srv.Name())
			}
			return nil
		})
	}

	defer s.log.Infow("service stopped")
	return g.Wait()
}

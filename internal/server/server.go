package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/ansedo/toptraffic/internal/logger"
	"github.com/ansedo/toptraffic/internal/router"
	"github.com/ansedo/toptraffic/internal/services/shutdowner"
)

type Server struct {
	http   http.Server
	logger *zap.Logger
}

func New(ctx context.Context, serverPort string, advDomains []string) *Server {
	s := &Server{
		http: http.Server{
			Addr:    serverPort,
			Handler: router.New(ctx, advDomains),
		},
		logger: logger.FromCtx(ctx),
	}
	s.addToShutdowner(ctx)
	return s
}

func (s *Server) Run(ctx context.Context) {
	go s.ListenAndServer(ctx)
}

func (s *Server) ListenAndServer(_ context.Context) {
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal(err.Error())
	}
}

func (s *Server) addToShutdowner(ctx context.Context) {
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		shutdown.AddCloser(func(ctx context.Context) error {
			if err := s.http.Shutdown(ctx); err != nil {
				return err
			}
			return nil
		})
	}
}

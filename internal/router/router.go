package router

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ansedo/toptraffic/internal/handlers"
)

func New(ctx context.Context, advDomains []string) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Compress(5),
	)

	r.Post("/placements/request", handlers.PlacementsRequest(ctx, advDomains))

	return r
}

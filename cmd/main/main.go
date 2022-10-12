package main

import (
	"context"

	"github.com/ansedo/toptraffic/internal/config"
	"github.com/ansedo/toptraffic/internal/logger"
	"github.com/ansedo/toptraffic/internal/server"
	"github.com/ansedo/toptraffic/internal/services/shutdowner"
)

func main() {
	// add logger to context
	ctx := logger.CtxWith(context.Background())

	// add shutdowner to context
	ctx = shutdowner.CtxWith(ctx)

	// construct config
	cfg := config.New(ctx)

	// construct and run server
	server.New(ctx, cfg.ServerPort, cfg.AdvDomains).Run(ctx)

	// waiting for graceful shutdown if it exists in context
	if shutdown := shutdowner.FromCtx(ctx); shutdown != nil {
		<-shutdown.ChShutdowned
	}
}

package main

import (
	"context"

	"github.com/jace-ys/roamd-world/backend/libraries/bedrock"
	"github.com/jace-ys/roamd-world/backend/libraries/config"
	"github.com/jace-ys/roamd-world/backend/libraries/logging"

	"github.com/jace-ys/roamd-world/backend/services/service.spots/handler"

	spotspb "github.com/jace-ys/roamd-world/backend/services/service.spots/proto/v1"
)

var cfg struct {
	config.StaticBase
}

func main() {
	config.MustParseStatic(&cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logging.MustNewLogger(cfg.LogLevel)
	defer logger.Sync()

	config.MustLoadDynamic(logger, cfg.Config)
	defer config.Stop()

	h := handler.New(logger)

	srv := bedrock.NewGRPCServer(logger, cfg.Port)
	srv.RegisterService(&spotspb.Spots_ServiceDesc, h)

	admin := bedrock.NewAdminServer(cfg.AdminPort)
	admin.Administer(srv)

	svc := bedrock.NewService(logger, srv, admin)
	if err := svc.Run(ctx); err != nil {
		logger.Fatalw("error running service", "error", err)
	}
}

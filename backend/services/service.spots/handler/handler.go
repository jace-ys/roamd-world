package handler

import (
	"context"

	"github.com/alexliesenfeld/health"
	"go.uber.org/zap"

	spotspb "github.com/jace-ys/roamd-world/backend/services/service.spots/proto/v1"
)

var _ spotspb.SpotsServer = (*Handler)(nil)

type Handler struct {
	spotspb.UnimplementedSpotsServer

	log *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (h *Handler) HealthChecks() []health.Check {
	return []health.Check{}
}

func (h *Handler) GetSpot(ctx context.Context, req *spotspb.GetSpotRequest) (*spotspb.GetSpotReply, error) {
	return &spotspb.GetSpotReply{}, nil
}

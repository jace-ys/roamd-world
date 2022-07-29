package handler

import (
	"context"

	"github.com/alexliesenfeld/health"
	"go.uber.org/zap"

	userspb "github.com/jace-ys/roamd-world/backend/services/service.users/proto/v1"
)

var _ userspb.UsersServer = (*Handler)(nil)

type Handler struct {
	userspb.UnimplementedUsersServer

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

func (h *Handler) GetUser(ctx context.Context, req *userspb.GetUserRequest) (*userspb.GetUserReply, error) {
	return &userspb.GetUserReply{}, nil
}

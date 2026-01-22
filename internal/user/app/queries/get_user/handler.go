package getuser

import (
	"context"
	"go-interview/internal/user/domain"
)

type Repository interface {
	GetUserByExternalID(ctx context.Context, externalID string) (*domain.User, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query GetUserByExternalIDQuery) (*GetUserByExternalIDResult, error) {
	user, err := h.repo.GetUserByExternalID(ctx, query.ExternalID)
	if err != nil {
		return nil, err
	}

	return &GetUserByExternalIDResult{
		ID:         user.ID.String(),
		CreatedAt:  user.CreatedAt.String(),
		ExternalID: user.ExternalID,
	}, nil
}

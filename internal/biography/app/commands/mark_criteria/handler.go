package markcriteria

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	MarkCriteria(ctx context.Context, ids ...uuid.UUID) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd MakrCriteriaCommand) error {
	ids := make([]uuid.UUID, 0, len(cmd.CriteriaIDs))
	for _, id := range cmd.CriteriaIDs {
		uuid, err := uuid.Parse(id)
		if err != nil {
			return err
		}
		ids = append(ids, uuid)
	}

	return h.repo.MarkCriteria(ctx, ids...)
}

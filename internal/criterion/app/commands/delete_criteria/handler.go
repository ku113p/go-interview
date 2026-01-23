package delete_criteria

import (
	"context"
	"go-interview/internal/criterion/domain"

	"github.com/google/uuid"
)

type DeleteCriteriaRepository interface {
	domain.CriteriaDeleter
}

type DeleteCriteriaHandler struct {
	repo DeleteCriteriaRepository
}

func NewDeleteCriteriaHandler(repo DeleteCriteriaRepository) *DeleteCriteriaHandler {
	return &DeleteCriteriaHandler{
		repo: repo,
	}
}

func (h *DeleteCriteriaHandler) Handle(ctx context.Context, cmd DeleteCriteriaCommand) error {
	toDeleteIDs := make([]uuid.UUID, 0, len(cmd.CriteriaIDs))
	for _, idStr := range cmd.CriteriaIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return err
		}
		toDeleteIDs = append(toDeleteIDs, id)
	}

	err := h.repo.DeleteCriteria(ctx, toDeleteIDs...)
	if err != nil {
		return err
	}

	return nil
}

package delete_criteria

import (
	"context"
	"errors"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type DeleteCriteriaRepository interface {
	domain.CriteriaDeleter
	domain.LifeAreaGetter
	domain.CriteriaNodeGetter
}

type DeleteCriteriaHandler struct {
	repo DeleteCriteriaRepository
}

func NewDeleteCriteriaHandler(repo DeleteCriteriaRepository) *DeleteCriteriaHandler {
	return &DeleteCriteriaHandler{
		repo: repo,
	}
}

func (h *DeleteCriteriaHandler) Handle(ctx context.Context, cmd DeleteCriteriaCommand) (*DeleteCriteriaResult, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	toDeleteIDs := make([]uuid.UUID, 0, len(cmd.CriteriaIDs))
	for _, idStr := range cmd.CriteriaIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		toDeleteIDs = append(toDeleteIDs, id)
	}

	criteriaNodeIDs, err := h.repo.GetCriteriaNodeIDs(ctx, toDeleteIDs...)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	checkedLifeAreas := make(map[uuid.UUID]*domain.LifeArea)
	for _, criterionID := range toDeleteIDs {
		lifeAreaID, ok := criteriaNodeIDs[criterionID]
		if !ok {
			return nil, domain.ErrNotFound
		}

		lifeArea, cached := checkedLifeAreas[lifeAreaID]
		if !cached {
			lifeArea, err = h.repo.GetLifeArea(ctx, lifeAreaID)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					return nil, domain.ErrNotFound
				}
				return nil, err
			}
			checkedLifeAreas[lifeAreaID] = lifeArea
		}

		if lifeArea.UserID != userID {
			return nil, domain.ErrForbidden
		}
	}

	err = h.repo.DeleteCriteria(ctx, toDeleteIDs...)
	if err != nil {
		return nil, err
	}

	return &DeleteCriteriaResult{}, nil
}

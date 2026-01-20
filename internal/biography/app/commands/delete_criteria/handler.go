package delete_criteria

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type DeleteCriteriaRepository interface {
	domain.CriteriaDeleter
	domain.LifeAreaGetter
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

	lifeAreas := make(map[uuid.UUID]*domain.LifeArea)
	for _, id := range toDeleteIDs {
		var lifeArea domain.LifeArea
		if lifeArea, ok := lifeAreas[id]; !ok {
			lifeArea, err = h.repo.GetLifeArea(ctx, id)
			if err != nil {
				return nil, err
			}
			lifeAreas[id] = lifeArea
		}
		if lifeArea.UserID != userID {
			return nil, fmt.Errorf("user %s is not owner of life area %s", userID, id)
		}
	}

	err = h.repo.DeleteCriteria(ctx, toDeleteIDs...)
	if err != nil {
		return nil, err
	}

	return &DeleteCriteriaResult{}, nil
}

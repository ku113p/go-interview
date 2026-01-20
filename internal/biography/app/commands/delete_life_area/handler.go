package delete_life_area

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type DeleteLifeAreaRepository interface {
	domain.LifeAreaDeleter
	domain.LifeAreaGetter
}

type DeleteLifeAreaHandler struct {
	repo DeleteLifeAreaRepository
}

func NewDeleteLifeAreaHandler(repo DeleteLifeAreaRepository) *DeleteLifeAreaHandler {
	return &DeleteLifeAreaHandler{
		repo: repo,
	}
}

func (h *DeleteLifeAreaHandler) Handle(ctx context.Context, cmd DeleteLifeAreaCommand) (*DeleteLifeAreaResult, error) {
	id, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	lifeArea, err := h.repo.GetLifeArea(ctx, id)
	if err != nil {
		return nil, err
	}

	if lifeArea.UserID != userID {
		return nil, fmt.Errorf("user %s is not owner of life area %s", userID, id)
	}

	err = h.repo.DeleteLifeArea(ctx, id)
	if err != nil {
		return nil, err
	}

	return &DeleteLifeAreaResult{}, nil
}

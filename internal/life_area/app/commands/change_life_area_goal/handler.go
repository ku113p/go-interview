package change_life_area_goal

import (
	"context"
	"errors"
	common "go-interview/internal/common/domain"
	"go-interview/internal/life_area/domain"

	"github.com/google/uuid"
)

type ChangeLifeAreaGoalRepository interface {
	domain.LifeAreaGetter
	domain.LifeAreaGoalChanger
}

type ChangeLifeAreaGoalHandler struct {
	repo ChangeLifeAreaGoalRepository
}

func NewChangeLifeAreaGoalHandler(repo ChangeLifeAreaGoalRepository) *ChangeLifeAreaGoalHandler {
	return &ChangeLifeAreaGoalHandler{
		repo: repo,
	}
}

func (h *ChangeLifeAreaGoalHandler) Handle(ctx context.Context, cmd ChangeLifeAreaGoalCommand) (*ChangeLifeAreaGoalResult, error) {
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
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	if lifeArea.UserID != userID {
		return nil, common.ErrForbidden
	}

	err = h.repo.ChangeGoal(ctx, lifeArea.ID, domain.NewGoal(cmd.Goal))
	if err != nil {
		return nil, err
	}

	return &ChangeLifeAreaGoalResult{}, nil
}

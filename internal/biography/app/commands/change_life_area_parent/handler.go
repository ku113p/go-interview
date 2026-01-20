package change_life_area_parent

import (
	"context"
	"fmt"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type ChangeLifeAreaParentRepository interface {
	domain.LifeAreaGetter
	domain.LifeAreaParentChanger
}

type ChangeLifeAreaParentHandler struct {
	repo ChangeLifeAreaParentRepository
}

func NewChangeLifeAreaParentHandler(repo ChangeLifeAreaParentRepository) *ChangeLifeAreaParentHandler {
	return &ChangeLifeAreaParentHandler{
		repo: repo,
	}
}

func (h *ChangeLifeAreaParentHandler) Handle(ctx context.Context, cmd ChangeLifeAreaParentCommand) (*ChangeLifeAreaParentResult, error) {
	id, err := uuid.Parse(cmd.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return nil, err
	}

	var parentID *uuid.UUID
	if cmd.ParentID != nil {
		parentPreID, err := uuid.Parse(*cmd.ParentID)
		if err != nil {
			return nil, err
		}
		parentID = &parentPreID
	}

	lifeArea, err := h.repo.GetLifeArea(ctx, id)
	if err != nil {
		return nil, err
	}

	if lifeArea.UserID != userID {
		return nil, fmt.Errorf("user %s is not owner of life area %s", userID, id)
	}
	err = h.repo.ChangeParentID(ctx, lifeArea.ID, parentID)
	if err != nil {
		return nil, err
	}

	return &ChangeLifeAreaParentResult{}, nil
}

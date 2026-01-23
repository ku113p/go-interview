package change_life_area_parent

import (
	"context"
	"errors"
	common "go-interview/internal/common/domain"
	"go-interview/internal/life_area/domain"

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
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	if lifeArea.UserID != userID {
		return nil, common.ErrForbidden
	}
	err = h.repo.ChangeParentID(ctx, lifeArea.ID, parentID)
	if err != nil {
		return nil, err
	}

	return &ChangeLifeAreaParentResult{}, nil
}

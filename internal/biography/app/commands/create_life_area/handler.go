package create_life_area

import (
	"context"
	"go-interview/internal/biography/domain"

	"github.com/google/uuid"
)

type CreateLifeAreaRepository interface {
	domain.LifeAreaCreator
	domain.LifeAreaGetter
}

type CreateLifeAreaHandler struct {
	repo  CreateLifeAreaRepository
	genID domain.IDGenerator
}

func NewCreateLifeAreaHandler(repo CreateLifeAreaRepository, genID domain.IDGenerator) *CreateLifeAreaHandler {
	return &CreateLifeAreaHandler{
		repo:  repo,
		genID: genID,
	}
}

func (h *CreateLifeAreaHandler) Handle(ctx context.Context, cmd CreateLifeAreaCommand) (*CreateLifeAreaResult, error) {
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

	id, err := h.genID.GenID()
	if err != nil {
		return nil, err
	}

	lifeArea := domain.NewLifeArea(
		id,
		userID,
		parentID,
		domain.NewTitle(cmd.Title),
		domain.NewGoal(cmd.Goal),
	)

	err = h.repo.CreateLifeArea(ctx, lifeArea)
	if err != nil {
		return nil, err
	}

	return &CreateLifeAreaResult{
		ID: lifeArea.ID.String(),
	}, nil
}

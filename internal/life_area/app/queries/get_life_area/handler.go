package get_life_area

import (
	"context"
	"errors"
	"go-interview/internal/life_area/domain"

	common "go-interview/internal/common/domain"

	"github.com/google/uuid"
)

type GetLifeAreaHandler struct {
	repo domain.LifeAreaGetter
}

func NewGetLifeAreaHandler(repo domain.LifeAreaGetter) *GetLifeAreaHandler {
	return &GetLifeAreaHandler{
		repo: repo,
	}
}

func (h *GetLifeAreaHandler) Handle(ctx context.Context, query GetLifeAreaQuery) (*GetLifeAreaResult, error) {
	queryID, err := uuid.Parse(query.ID)
	if err != nil {
		return nil, err
	}

	la, err := h.repo.GetLifeArea(ctx, queryID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	var parentID *string
	if la.ParentID != nil {
		idStr := la.ParentID.String()
		parentID = &idStr
	}

	return &GetLifeAreaResult{
		ID:        la.ID.String(),
		Title:     la.Title.String(),
		Goal:      la.Goal.String(),
		ParentID:  parentID,
		CreatedAt: la.CreatedAt.String(),
		UpdatedAt: la.UpdatedAt.String(),
	}, nil
}

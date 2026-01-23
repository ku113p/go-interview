package list_live_area

import (
	"context"
	"go-interview/internal/life_area/domain"

	"github.com/google/uuid"
)

type ListLifeAreaHandler struct {
	repo  domain.LifeAreaLister
	genID domain.IDGenerator
}

func NewListLifeAreaHandler(repo domain.LifeAreaLister, genID domain.IDGenerator) *ListLifeAreaHandler {
	return &ListLifeAreaHandler{
		repo:  repo,
		genID: genID,
	}
}

func (h *ListLifeAreaHandler) Handle(ctx context.Context, query ListLifeAreaQuery) (*ListLifeAreaResult, error) {
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, err
	}

	las, err := h.repo.ListLifeAreas(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*LifeArea, 0, len(las))
	for _, la := range las {
		var parentID *string
		if la.ParentID != nil {
			idStr := la.ParentID.String()
			parentID = &idStr
		}
		result = append(result, &LifeArea{
			ID:        la.ID.String(),
			Title:     la.Title.String(),
			Goal:      la.Goal.String(),
			ParentID:  parentID,
			CreatedAt: la.CreatedAt.String(),
			UpdatedAt: la.UpdatedAt.String(),
		})
	}
	return &ListLifeAreaResult{
		Items: result,
	}, nil
}

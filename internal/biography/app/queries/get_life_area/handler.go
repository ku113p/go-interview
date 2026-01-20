package get_life_area

import (
	"context"
	"go-interview/internal/biography/domain"

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
		return nil, err
	}

	var parentID *string
	if la.ParentID != nil {
		idStr := la.ParentID.String()
		parentID = &idStr
	}

	criteria := make([]*Criterion, 0, len(la.Criteria))
	for _, c := range la.Criteria {
		criteria = append(criteria, &Criterion{
			ID:          c.ID.String(),
			Description: c.Description.String(),
			IsCompleted: c.IsCompleted,
		})
	}

	return &GetLifeAreaResult{
		ID:        la.ID.String(),
		Title:     la.Title.String(),
		Goal:      la.Goal.String(),
		ParentID:  parentID,
		Criteria:  criteria,
		CreatedAt: la.CreatedAt.String(),
		UpdatedAt: la.UpdatedAt.String(),
	}, nil
}

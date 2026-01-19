package get

import (
	"context"
	"go-interview/internal/biography/domain"
)

type Handler struct {
	repo domain.LifeAreaRepository
}

func NewHandler(repo domain.LifeAreaRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) Handle(ctx context.Context, query Query) (Result, error) {
	la, err := h.repo.GetLifeArea(ctx, query.ID)
	if err != nil {
		return Result{}, err
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

	return Result{
		ID:        la.ID.String(),
		Title:     la.Title.String(),
		Goal:      la.Goal.String(),
		ParentID:  parentID,
		Criteria:  criteria,
		CreatedAt: la.CreatedAt.String(),
		UpdatedAt: la.UpdatedAt.String(),
	}, nil
}
